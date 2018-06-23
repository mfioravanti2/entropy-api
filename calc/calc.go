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

type Attribute struct {
	Mnemonic string
	Locale   string
}

type Attributes []Attribute

type Person struct {
	H float64
	Attributes mapset.Set
	Heuristics []string
	Errors []error
}

func Calc( ctx context.Context, r *request.Request, formatId string ) (response.Response, error) {
	var err error
	var errors []error
	var score response.Response
	var attributes mapset.Set
	var h_t, h_total float64
	var nationality string

	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "Calc", "calc" )
	} else {
		ctx = logging.WithFuncId( ctx, "Calc", "calc" )
	}

	logger := logging.Logger( ctx )

	logger.Info("preparing to score request",
	)

	nations := make( map[string]*source.Model )
	h_total = 0.0
	attributes = mapset.NewSet()

	nationality = strings.ToLower( r.Locale )
	var m *source.Model
	if m, err = data.GetModel(nationality); err != nil {
		score.Errors = new(response.Errors)
		score.Errors.Messages = append( score.Errors.Messages, err.Error() )
		return score, err
	}
	nations[nationality] = m
	t := source.Threshold{ Locale: r.Locale, Threshold: m.Threshold, K: m.K}

	for _, p := range r.People {
		nation := strings.ToLower( p.Nationality )

		if _, ok := nations[nation]; !ok {
			var n *source.Model

			if n, err = data.GetModel(nation); err != nil {
				score.Errors = new(response.Errors)
				score.Errors.Messages = append( score.Errors.Messages, err.Error() )
				return score, err
			}

			if n.Threshold < t.Threshold {
				t = source.Threshold{ Locale: r.Locale, Threshold: n.Threshold, K: n.K}
			}

			nations[nation] = n
		}

		person := calcPerson( ctx, p, nations[nation], formatId )
		attributes = attributes.Union( ConvertSet( person.Attributes, nation) )
		h_total += person.H

		if len(errors) > 0 {
			if score.Errors == nil {
				score.Errors = new(response.Errors)
			}

			for _, e := range person.Errors {
				score.Errors.Messages = append( score.Errors.Messages, e.Error() )
			}
		}
	}

	h_t = t.Threshold

	score.Data = new(response.Data)
	score.Data.Pii = h_total >= h_t
	score.Data.Locale = t.Locale
	score.Data.Score = h_total
	score.Data.RunDate = time.Now().UTC()
	score.Data.ApiVersion = sys.VERSION

	for val := range attributes.Iterator().C {
		if a, ok := val.(Attribute); ok {
			s, err := source.GetScore( nations[a.Locale], a.Mnemonic, formatId)
			if err == nil {
				r := response.Attribute{Mnemonic: a.Mnemonic, Locale: strings.ToUpper(a.Locale), Format: formatId, Score: s}
				score.Data.Attributes = append(score.Data.Attributes, r)
			}
		}
	}

	return score, nil
}

func ArrayToSet( a []string ) mapset.Set {
	m := mapset.NewSet()

	for _, s := range a {
		m.Add( s )
	}

	return m
}

func ConvertSet( m mapset.Set, locale string ) mapset.Set {
	f := mapset.NewSet()

	for val := range m.Iterator().C {
		if str, ok := val.(string); ok {
			f.Add(Attribute{ Mnemonic: str, Locale: locale })
		}
	}

	return f
}

func SetToArray( m mapset.Set ) []string {
	var a []string

	for val := range m.Iterator().C {
		if str, ok := val.(string); ok {
			a = append(a, str)
		}
	}

	return a
}

func containsAll( m mapset.Set, s []string) bool {
	for _, e := range s {
		if !m.Contains(e) {
			return false
		}
	}

	return true
}

func calcPerson( ctx context.Context, p request.Person, s *source.Model, formatId string ) Person {
	var changed bool = false
	var loops, changes int = 0, 0
	var person Person
	person.Attributes = mapset.NewSet()
	person.H = 0.0

	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "calcPerson", "calc" )
	} else {
		ctx = logging.WithFuncId( ctx, "calcPerson", "calc" )
	}

	logger := logging.Logger( ctx )

	logger.Info("preparing to score individual",
		zap.String( "personId", p.PersonID ),
	)

	for _, a := range p.Attributes {
		person.Attributes.Add(a.Mnemonic)
	}

	for {
		changed = false

		for _, h := range s.Heuristics {
			logger.Debug("checking heuristic attribute set",
				zap.String( "personId", p.PersonID ),
				zap.String( "heuristicId", h.Id ),
				zap.Int( "loopId", loops ),
			)

			if containsAll(person.Attributes, h.Match) {
				logger.Info("heuristic attribute set matched",
					zap.String( "personId", p.PersonID ),
					zap.String( "heuristicId", h.Id ),
					zap.Int( "loopId", loops ),
				)

				if len(h.Remove) > 0 {
					logger.Info("removing attribute set based on heuristic match",
						zap.String( "personId", p.PersonID ),
						zap.String( "heuristicId", h.Id ),
						zap.Int( "loopId", loops ),
					)

					r_s := ArrayToSet( h.Remove )
					person.Attributes = person.Attributes.Difference( r_s )
					changed = true
				}

				if len(h.Insert) > 0 {
					logger.Info("inserting attribute set based on heuristic match",
						zap.String( "personId", p.PersonID ),
						zap.String( "heuristicId", h.Id ),
						zap.Int( "loopId", loops ),
					)

					r_i := ArrayToSet( h.Insert )
					person.Attributes = person.Attributes.Union( r_i )
					changed = true
				}

				changes += 1
			}
		}

		if !changed {
			logger.Info("heuristics comparisons completed, no more changes",
				zap.String( "personId", p.PersonID ),
				zap.Int( "loopId", loops ),
				zap.Int( "changes", changes ),
			)

			break
		}

		loops += 1
	}

	for val := range person.Attributes.Iterator().C {
		if attributeId, ok := val.(string); ok {
			h_i, err := source.GetScore( s, attributeId, formatId )
			if err == nil {
				logger.Info("scoring final attribute set",
					zap.String( "personId", p.PersonID ),
					zap.String( "formatId", formatId ),
					zap.String( "attributeId", attributeId ),
					zap.Float64( "score", h_i ),
				)

				person.H += h_i
			} else {
				logger.Info("scoring final attribute set",
					zap.String( "personId", p.PersonID ),
					zap.String( "formatId", formatId ),
					zap.String( "attributeId", attributeId ),
					zap.String( "status", "error" ),
					zap.String( "error", err.Error() ),
				)

				person.Errors = append( person.Errors, err )
			}
		}
	}

	return person
}
