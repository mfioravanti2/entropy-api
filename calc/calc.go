package calc

import (
	"time"
	"strings"
	"context"

	"go.uber.org/zap"
	"github.com/deckarep/golang-set"

	"github.com/mfioravanti2/entropy-api/model/response"
	"github.com/mfioravanti2/entropy-api/model/request"
	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model/source"
	"github.com/mfioravanti2/entropy-api/api/sys"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type Person struct {
	H float64
	Heuristics response.Heuristics
	Tags map[string]mapset.Set
	Errors []error
}

// Calculate the Entropy scores for an attribute set
func Calc( ctx context.Context, r *request.Request, formatId string, useReductions bool ) (response.Response, error) {
	var err error
	var score response.Response
	var h_t, h_total float64
	var nationality string

	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "Calc", "calc" )
	} else {
		ctx = logging.WithFuncId( ctx, "Calc", "calc" )
	}

	logger := logging.Logger( ctx )

	logger.Info("scoring of individual attribute sets has started" )

	score.Data = new(response.Data)
	nations := make( map[string]*source.Model )
	h_total = 0.0

	nationality = strings.ToLower( r.Locale )
	var m *source.Model
	if m, err = data.GetModel(nationality); err != nil {
		score.Errors = new(response.Errors)
		score.Errors.Messages = append( score.Errors.Messages, err.Error() )
		return score, err
	}
	nations[nationality] = m
	t := source.Threshold{ Locale: r.Locale, Threshold: m.Threshold, K: m.K}

	// Inspect the set of attributes associated with an individual person
	for _, p := range r.People {
		nation := strings.ToLower( p.Nationality )

		if _, ok := nations[nation]; !ok {
			var n *source.Model

			// retrieve the country model associated with the person's nationality
			if n, err = data.GetModel(nation); err != nil {
				score.Errors = new(response.Errors)
				score.Errors.Messages = append( score.Errors.Messages, err.Error() )
				return score, err
			}

			// if the peron's nationality's country model has a lower threshold,
			// use that model for the over all threshold determination
			// NOTE: t - K is not factored into the determination
			if n.Threshold < t.Threshold {
				t = source.Threshold{ Locale: r.Locale, Threshold: n.Threshold, K: n.K}
			}

			nations[nation] = n
		}

		var person Person
		if useReductions {
			// score the attribute set of an individual, with heuristic reductions
			person = calcReducePerson( ctx, p, nations[nation], formatId )
		} else {
			// score the attribute set of an individual without applying heuristic reductions
			person = calcPerson( ctx, p, nations[nation], formatId )
		}

		// convert the Person scoring object into a response object
		var cPerson response.Person
		cPerson.Id = p.PersonID
		cPerson.Nationality = p.Nationality
		cPerson.Score = person.H

		// if any heuristics have been applied, add those heuristics to the response list
		if len(person.Heuristics) > 0 {
			cPerson.Heuristics = new(response.Heuristics)
			*(cPerson.Heuristics) = append(*(cPerson.Heuristics), person.Heuristics...)
		}

		for tagId, t := range person.Tags {
			for val := range t.Iterator().C {
				if attributeId, ok := val.(string); ok {
					// collect the individual attribute scores and add those to the response
					s, err := nations[nation].Score( attributeId, formatId )
					if err == nil {
						r := response.Attribute{Mnemonic: attributeId, Locale: strings.ToUpper(nation), Tag: tagId, Format: formatId, Score: s}
						cPerson.Attributes = append( cPerson.Attributes, r )
					}
				}
			}
		}

		// Add the individual's entropy score to the overall request's entropy score
		h_total += cPerson.Score

		// If any errors were encountered during the scoring process, add those
		// errors to the overall response's error list
		if len(person.Errors) > 0 {
			if score.Errors == nil {
				score.Errors = new(response.Errors)
			}

			for _, e := range person.Errors {
				score.Errors.Messages = append( score.Errors.Messages, e.Error() )
			}
		}

		score.Data.People = append( score.Data.People, cPerson )
	}

	logger.Info("scoring of individual attribute sets are complete" )

	h_t = t.Threshold

	score.Data.Pii = h_total >= h_t
	score.Data.Locale = t.Locale
	score.Data.Score = h_total
	score.Data.RunDate = time.Now().UTC()
	score.Data.ApiVersion = sys.VERSION

	return score, nil
}

// Convert an array to a set
func ArrayToSet( a []string ) mapset.Set {
	m := mapset.NewSet()

	for _, s := range a {
		m.Add( s )
	}

	return m
}

// Determine if a set contains all of the attributes from another set
func containsAll( m mapset.Set, s []string) bool {
	for _, e := range s {
		if !m.Contains(e) {
			return false
		}
	}

	return true
}

// Calculate an individual person's attribute set entropy score,
// allow the scores to be modified by applying the heuristics rule sets
func calcReducePerson( ctx context.Context, p request.Person, s *source.Model, formatId string ) Person {
	var changed bool = false
	var loops, changes int
	var person Person
	person.Tags = make( map[string]mapset.Set )
	person.H = 0.0

	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "calcReducePerson", "calc" )
	} else {
		ctx = logging.WithFuncId( ctx, "calcReducePerson", "calc" )
	}

	logger := logging.Logger( ctx )

	logger.Info("scoring of an individual's attribute set is starting",
		zap.String( "personId", p.PersonID ),
	)

	// collect a set of tags
	for _, a := range p.Attributes {
		if _, ok := person.Tags[a.Tag]; !ok {
			person.Tags[a.Tag] = mapset.NewSet()

			logger.Debug("registering new tag in attribute set",
				zap.String( "personId", p.PersonID ),
				zap.String( "tagId", a.Tag ),
			)
		}

		person.Tags[a.Tag].Add(a.Mnemonic)
	}

	// score all attributes associated with a specific tag
	changes = 0
	for tagId, t := range person.Tags {
		loops = 0

		logger.Info("scoring attribute set with tag",
			zap.String( "personId", p.PersonID ),
			zap.String( "tagId", tagId ),
		)

		// Loop through the set of attributes and terminate the application
		// of the reduction heuristics when no more changes to the attribute
		// set have occurred
		for {
			changed = false

			for _, h := range s.Heuristics {
				logger.Debug("checking heuristic attribute set",
					zap.String( "personId", p.PersonID ),
					zap.String( "heuristicId", h.Id ),
					zap.String( "tagId", tagId ),
					zap.Int( "loopId", loops ),
				)

				// determine if the attribute set contains all of the
				// heuristics' attributes
				if containsAll( t, h.Match) {
					logger.Debug("heuristic attribute set matched",
						zap.String( "personId", p.PersonID ),
						zap.String( "heuristicId", h.Id ),
						zap.String( "tagId", tagId ),
						zap.Int( "loopId", loops ),
					)

					// record the heuristics as being applied to the attribute set
					person.Heuristics = append( person.Heuristics, h.Id)

					// if a heuristic is applied, remove the attributes specified by
					// the matched heuristic
					if len(h.Remove) > 0 {
						logger.Debug("removing attribute set based on heuristic match",
							zap.String( "personId", p.PersonID ),
							zap.String( "heuristicId", h.Id ),
							zap.String( "tagId", tagId ),
							zap.Int( "loopId", loops ),
						)

						r_s := ArrayToSet( h.Remove )
						t = t.Difference( r_s )
						changed = true
					}

					// if a heuristic is applied, add all of the attributes specified by
					// the matched heuristic
					if len(h.Insert) > 0 {
						logger.Debug("inserting attribute set based on heuristic match",
							zap.String( "personId", p.PersonID ),
							zap.String( "heuristicId", h.Id ),
							zap.String( "tagId", tagId ),
							zap.Int( "loopId", loops ),
						)

						r_i := ArrayToSet( h.Insert )
						t = t.Union( r_i )
						changed = true
					}

					// record that a heuristic rule set has been matched and applied
					changes += 1
				}
			}

			if !changed {
				// no heuristics have been applied, terminate checks
				logger.Debug("heuristics comparisons completed, no more changes",
					zap.String( "personId", p.PersonID ),
					zap.String( "tagId", tagId ),
					zap.Int( "loopId", loops ),
					zap.Int( "changes", changes ),
				)

				break
			}

			loops += 1
		}

		person.Tags[tagId] = t

		for val := range t.Iterator().C {
			if attributeId, ok := val.(string); ok {

				// annotate the individual scores for each attribute
				h_i, err := s.Score( attributeId, formatId )
				if err == nil {
					logger.Info("scoring final attribute set",
						zap.String( "personId", p.PersonID ),
						zap.String( "tagId", tagId ),
						zap.String( "formatId", formatId ),
						zap.String( "attributeId", attributeId ),
						zap.Float64( "score", h_i ),
					)

					person.H += h_i
				} else {
					logger.Error("scoring final attribute set",
						zap.String( "personId", p.PersonID ),
						zap.String( "tagId", tagId ),
						zap.String( "formatId", formatId ),
						zap.String( "attributeId", attributeId ),
						zap.String( "status", "error" ),
						zap.String( "error", err.Error() ),
					)

					person.Errors = append( person.Errors, err )
				}
			}
		}
	}

	logger.Info("scoring of an individual's attribute set is complete",
		zap.String( "personId", p.PersonID ),
	)

	return person
}

// Calculate an individual person's attribute set entropy score, but DO NOT
// allow the scores to be modified by applying the heuristics rule sets
func calcPerson( ctx context.Context, p request.Person, s *source.Model, formatId string ) Person {
	var person Person
	person.Tags = make( map[string]mapset.Set )
	person.H = 0.0

	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "calcPerson", "calc" )
	} else {
		ctx = logging.WithFuncId( ctx, "calcPerson", "calc" )
	}

	logger := logging.Logger( ctx )

	logger.Info("scoring of an individual's attribute set is starting",
		zap.String( "personId", p.PersonID ),
	)

	// collect a set of tags
	for _, a := range p.Attributes {
		if _, ok := person.Tags[a.Tag]; !ok {
			person.Tags[a.Tag] = mapset.NewSet()

			logger.Debug("registering new tag in attribute set",
				zap.String( "personId", p.PersonID ),
				zap.String( "tagId", a.Tag ),
			)
		}

		person.Tags[a.Tag].Add(a.Mnemonic)
	}

	// score all attributes associated with a specific tag
	for tagId, t := range person.Tags {

		logger.Info("scoring attribute set with tag",
			zap.String( "personId", p.PersonID ),
			zap.String( "tagId", tagId ),
		)

		for val := range t.Iterator().C {
			if attributeId, ok := val.(string); ok {

				// annotate the individual scores for each attribute
				h_i, err := s.Score( attributeId, formatId )
				if err == nil {
					logger.Info("scoring final attribute set",
						zap.String( "personId", p.PersonID ),
						zap.String( "tagId", tagId ),
						zap.String( "formatId", formatId ),
						zap.String( "attributeId", attributeId ),
						zap.Float64( "score", h_i ),
					)

					person.H += h_i
				} else {
					logger.Error("scoring final attribute set",
						zap.String( "personId", p.PersonID ),
						zap.String( "tagId", tagId ),
						zap.String( "formatId", formatId ),
						zap.String( "attributeId", attributeId ),
						zap.String( "status", "error" ),
						zap.String( "error", err.Error() ),
					)

					person.Errors = append( person.Errors, err )
				}
			}
		}
	}

	logger.Info("scoring of an individual's attribute set is complete",
		zap.String( "personId", p.PersonID ),
	)

	return person
}
