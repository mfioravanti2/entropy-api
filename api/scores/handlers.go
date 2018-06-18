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
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func AddHandlers(r model.Routes) model.Routes {
	logger := logging.Logger(nil)
	logger.Info("AddHandlers(scores)")

	r = append( r, model.Route{"ScoreCalc", "POST", "/v1/scores", Calc} )
	r = append( r, model.Route{"ScoreCalcFormat", "POST", "/v1/scores/format/{formatId}", CalcFormat} )

	return r
}

func Calc(w http.ResponseWriter, r *http.Request) {
	score( w, r, "mean")
}

func CalcFormat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	formatId := strings.ToLower(vars["formatId"])

	score( w, r, formatId )
}

func score(w http.ResponseWriter, r *http.Request, formatId string) {
	var entropy request.Request

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 50 * 1024 ))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &entropy); err != nil {
		var s = "invalid request object, expected json format"
		handleError( w, r, http.StatusUnprocessableEntity, s )
		return
	}

	var score response.Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if score, err = calc.Calc( &entropy, formatId ); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(score); err != nil {
		panic(err)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	var score response.Response

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader( statusCode )

	score.Errors = new(response.Errors)
	score.Errors.Messages = append( score.Errors.Messages, msg )
	if err := json.NewEncoder(w).Encode(score); err != nil {
		panic(err)
	}
}

