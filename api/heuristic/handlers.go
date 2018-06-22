package heuristic

import (
	"context"
	"net/http"
	"encoding/json"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "heuristic" )

	logger := logging.Logger( ctx )

	logger.Info("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/heuristics" ) )
	r = append( r, model.Route{"HeuristicList", "GET", "/v1/countries/{countryId}/heuristics", List})

	logger.Info("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/heuristics/{heuristicId}" ) )
	r = append( r, model.Route{"HeuristicDetails", "GET", "/v1/countries/{countryId}/heuristics/{heuristicId}", Detail})

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])

	var err error
	var heuristics []string
	heuristics, err = data.GetHeuristics(countryId)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if len(heuristics) > 0 {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(heuristics); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}

func Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])
	heuristicId := strings.ToLower(vars["heuristicId"])

	heuristic, err := data.GetHeuristic(countryId, heuristicId)

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if err == nil {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(heuristic); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}
