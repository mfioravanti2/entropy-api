package calc

import (
	"fmt"
	"time"

	"github.com/mfioravanti2/entropy-api/model/response"
	"github.com/mfioravanti2/entropy-api/model/request"
	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model/source"
	"strings"
)

func Calc( r *request.Request ) (response.Response, error) {
	var err error
	var score response.Response
	var a_p []string
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

		a_p, h_p = calcPerson( p, nations[nation] )
		fmt.Println("set:", a_p)
		h_total += h_p

	}

	score.Pii = h_total >= h_t
	score.Score = h_total
	score.RunDate = time.Now()

	return score, nil
}

func calcPerson( p request.Person, s source.Model ) ([]string, float64) {
	var attributes []string
	var h_p float64 = 0.0

	return attributes, h_p
}
