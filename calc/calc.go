package calc

import (
	"time"
	"strings"

	"github.com/deckarep/golang-set"

	"github.com/mfioravanti2/entropy-api/model/response"
	"github.com/mfioravanti2/entropy-api/model/request"
	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model/source"
	"github.com/mfioravanti2/entropy-api/api/sys"
)

type Attribute struct {
	Name string
	Locale string
}

func Calc( r *request.Request, formatId string ) (response.Response, error) {
	var err error
	var score response.Response
	var attributes, a_p mapset.Set
	var h_t, h_p, h_total float64
	var nationality string

	nations := make( map[string]source.Model )
	h_total = 0.0
	attributes = mapset.NewSet()

	nationality = strings.ToLower( r.Locale )
	if nations[nationality], err = data.GetModel(nationality); err != nil {
		score.Errors = new(response.Errors)
		score.Errors.Messages = append( score.Errors.Messages, err.Error() )
		return score, err
	}
	t := source.Threshold{ Locale: r.Locale, Threshold: nations[r.Locale].Threshold, K: nations[r.Locale].K}

	for _, p := range r.People {
		nation := strings.ToLower( p.Nationality )

		if _, ok := nations[nation]; !ok {
			if nations[nation], err = data.GetModel(nation); err != nil {
				score.Errors = new(response.Errors)
				score.Errors.Messages = append( score.Errors.Messages, err.Error() )
				return score, err
			}

			if nations[nation].Threshold < t.Threshold {
				t = source.Threshold{ Locale: r.Locale, Threshold: nations[nation].Threshold, K: nations[nation].K}
			}
		}

		a_p, h_p = calcPerson( p, nations[nation], formatId )
		attributes = attributes.Union( ConvertSet( a_p, nation) )
		h_total += h_p
	}

	h_t = t.Threshold

	score.Data = new(response.Data)
	score.Data.Pii = h_total >= h_t
	score.Data.Locale = t.Locale
	score.Data.Score = h_total
	score.Data.RunDate = time.Now()
	score.Data.ApiVersion = sys.VERSION

	for val := range attributes.Iterator().C {
		if a, ok := val.(Attribute); ok {
			s, err := source.GetScore(nations[a.Locale], a.Name, formatId)
			if err == nil {
				r := response.Attribute{Name: a.Name, Locale: strings.ToUpper(a.Locale), Format: formatId, Score: s}
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
			f.Add(Attribute{ Name: str, Locale: locale })
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

func calcPerson( p request.Person, s source.Model, formatId string ) (mapset.Set, float64) {
	var h_p float64 = 0.0
	var changed bool = false
	a_p := mapset.NewSet()

	for _, a := range p.Attributes {
		a_p.Add(a.Name)
	}

	for {
		changed = false

		for _, h := range s.Heuristics {
			if containsAll(a_p, h.Match) {
				if len(h.Remove) > 0 {
					r_s := ArrayToSet( h.Remove )
					a_p = a_p.Difference( r_s )
					changed = true
				}

				if len(h.Insert) > 0 {
					r_i := ArrayToSet( h.Insert )
					a_p = a_p.Union( r_i )
					changed = true
				}
			}
		}

		if !changed {
			break
		}
	}

	for val := range a_p.Iterator().C {
		if str, ok := val.(string); ok {
			h_i, err := source.GetScore(s, str, formatId)
			if err == nil {
				h_p += h_i
			}
		}
	}

	return a_p, h_p
}
