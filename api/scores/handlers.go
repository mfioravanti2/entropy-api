package scores

import (
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/request"
	"github.com/mfioravanti2/entropy-api/model/response"
	"github.com/mfioravanti2/entropy-api/calc"
)

func AddHandlers(r model.Routes) model.Routes {
	r = append( r, model.Route{"ScoreCalc", "POST", "/v1/scores", Calc} )
	r = append( r, model.Route{"ScoreCalcFormat", "POST", "/v1/scores/format/{formatId}", CalcFormat} )

	return r
}

func Calc(w http.ResponseWriter, r *http.Request) {
	var entropy request.Request

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 50 * 1024))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &entropy); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader( http.StatusUnprocessableEntity )
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	var score response.Response
	if score, err = calc.Calc( &entropy, "raw" ); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(score); err != nil {
		panic(err)
	}
}

func CalcFormat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	formatId := strings.ToLower(vars["formatId"])

	var entropy request.Request

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 50 * 1024))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &entropy); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader( http.StatusUnprocessableEntity )
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	var score response.Response
	if score, err = calc.Calc( &entropy, formatId ); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(score); err != nil {
		panic(err)
	}
}
