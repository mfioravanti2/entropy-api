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

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/heuristics" ) )
	r = append( r, model.Route{"HeuristicList", "GET", "/v1/countries/{countryId}/heuristics", List, nil})

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/heuristics/{heuristicId}" ) )
	r = append( r, model.Route{"HeuristicDetails", "GET", "/v1/countries/{countryId}/heuristics/{heuristicId}", Detail, nil})

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	if ok, _ := model.ValidateCountryCode(countryId); !ok {
		logger.Error( "validating country code",
			zap.String("countryId", strings.ToUpper(countryId)),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid country code" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	logger.Info( "retrieving heuristics from country model",
		zap.String("countryId", strings.ToUpper(countryId)),
	)

	var err error
	var heuristics []string
	heuristics, err = data.GetHeuristics(countryId)
	if err != nil {
		logger.Error( "retrieving heuristics from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if len(heuristics) > 0 {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(heuristics); err != nil {
			logger.Error( "encoding heuristics",
				zap.String("countryId", strings.ToUpper(countryId) ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )

		logger.Info( "retrieving heuristics from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "status", "ok" ),
			zap.String("error ", "no heuristics found" ),
		)
	}
}

func Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])
	heuristicId := strings.ToLower(vars["heuristicId"])

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	if ok, _ := model.ValidateCountryCode(countryId); !ok {
		logger.Error( "validating country code",
			zap.String("countryId", strings.ToUpper(countryId)),
			zap.String( "heuristicId", strings.ToLower(heuristicId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid country code" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	if ok, _ := model.ValidateHeuristic(heuristicId); !ok {
		logger.Error( "validating heuristic identifier",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "heuristicId", strings.ToLower(heuristicId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid heuristic identifier" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	logger.Debug( "retrieving heuristic from country model",
		zap.String("countryId", strings.ToUpper(countryId) ),
		zap.String( "heuristicId", strings.ToLower(heuristicId) ),
	)

	heuristic, err := data.GetHeuristic(countryId, heuristicId)

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if err == nil {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(heuristic); err != nil {
			logger.Error( "encoding heuristic",
				zap.String("countryId", strings.ToUpper(countryId) ),
				zap.String( "heuristicId", strings.ToLower(heuristicId) ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )

		logger.Error( "retrieving heuristic from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "heuristicId", strings.ToLower(heuristicId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}
}
