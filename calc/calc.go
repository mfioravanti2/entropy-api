package calc

import (
	"time"
	"strings"

	"github.com/deckarep/golang-set"

	"github.com/mfioravanti2/entropy-api/model/response"
	"github.com/mfioravanti2/entropy-api/model/request"
	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model/source"
)

func Calc( r *request.Request, formatId string ) (response.Response, error) {
	var err error
	var score response.Response
	var attributes, a_p mapset.Set
	var h_t, h_p, h_total float64
	var nationality string

	nations := make( map[string]source.Model )
	h_total = 0.0

	nationality = strings.ToLower( r.Locale )
	if nations[nationality], err = data.GetModel(nationality); err != nil {
		panic(err)
	}
	h_t = nations[r.Locale].Threshold

	for _, p := range r.People {
		nation := strings.ToLower( p.Nationality )

		if _, ok := nations[nation]; !ok {
			if nations[nation], err = data.GetModel(nation); err != nil {
				panic(err)
			}

			if nations[nation].Threshold < h_t {
				h_t = nations[nation].Threshold
			}
		}

		a_p, h_p = calcPerson( p, nations[nation], formatId )
		attributes = attributes.Union( a_p )
		h_total += h_p
	}

	score.Pii = h_total >= h_t
	score.Score = h_total
	score.RunDate = time.Now()

	for val := range attributes.Iterator().C {
		if str, ok := val.(string); ok {
			score.Attributes = append( score.Attributes, response.Attribute{ str, formatId, 0.0})
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

func SetToArray( m mapset.Set ) []string {
	var a []string

	for val := range m.Iterator().C {
		if str, ok := val.(string); ok {
			a = append(a, str)
		}
	}

	return a
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
			h_s := ArrayToSet(h.Match)

			if a_p.Contains( h_s ) {
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
