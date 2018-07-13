package scores

import (
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"context"
	"time"

	"go.uber.org/zap"
	"github.com/gorilla/mux"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/request"
	"github.com/mfioravanti2/entropy-api/model/response"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/calc"
	"github.com/mfioravanti2/entropy-api/data/scoringdb"
)

const (
	// valid values are naive, mean, or rare
	DEFAULT_SCORING = "mean"
	// valid values are detailed or summary
	DEFAULT_MODE = "detailed"
	// valide values are exclude or include
	DEFAULT_REDUCTIONS = "include"
)

func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "scores" )

	logger := logging.Logger( ctx )

	p := []string{"format", "{formatId}"}

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"DetailedScoring", "POST", "/v1/scores", CalcOptions, p} )

	p = []string{"reductions", "{useReductions}"}

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"DetailedScoring", "POST", "/v1/scores", CalcOptions, p} )

	p = []string{"mode", "{modeId}"}

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"DetailedScoring", "POST", "/v1/scores", CalcOptions, p} )

	p = []string{"format", "{formatId}", "reductions", "{useReductions}"}

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"DetailedScoring", "POST", "/v1/scores", CalcOptions, p} )

	p = []string{"format", "{formatId}", "mode", "{modeId}"}

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"DetailedScoring", "POST", "/v1/scores", CalcOptions, p} )

	p = []string{"format", "{formatId}", "mode", "{modeId}", "reductions", "{useReductions}"}

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"DetailedScoring", "POST", "/v1/scores", CalcOptions, p} )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"DefaultScoring", "POST", "/v1/scores", CalcDefaults, nil} )

	return r
}

func CalcDefaults(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	logger.Info( "preparing to score request, with default formatId",
		zap.String("formatId", strings.ToLower(DEFAULT_SCORING) ),
	)

	score( w, r, DEFAULT_MODE, DEFAULT_SCORING, true )
}

func CalcOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	formatId, ok := vars["formatId"]
	if !ok || formatId == "" {
		formatId = DEFAULT_SCORING
	}

	modeId, ok := vars["modeId"]
	if !ok || modeId == "" {
		modeId = DEFAULT_MODE
	}

	useReductions, ok := vars["useReductions"]
	if !ok || useReductions == "" {
		useReductions = DEFAULT_REDUCTIONS
	}

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	logger.Info( "preparing to score request, with request specified formatId",
		zap.String("formatId", strings.ToLower(formatId) ),
		zap.String("modeId", strings.ToLower(modeId) ),
		zap.Bool("useReductions", useReductions == "include" ),
	)

	score( w, r, modeId, formatId, useReductions == "include" )
}

func score(w http.ResponseWriter, r *http.Request, modeId string, formatId string, useReductions bool ) {
	var entropy request.Request

	reqCtx := logging.WithFuncId( r.Context(), "score", "scores" )
	logger := logging.Logger(reqCtx)

	if ok, _ := model.ValidateFormat(formatId); !ok {
		logger.Error( "validating format identifier",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid format identifier" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	logger.Info( "score request",
		zap.String("formatId", strings.ToLower(formatId) ),
		zap.String("modeId", strings.ToLower(modeId) ),
		zap.Bool("useReductions", useReductions ),
	)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 50 * 1024 ))
	if err != nil {
		logger.Error( "unable to read request body",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}
	if err := r.Body.Close(); err != nil {
		logger.Error( "unable to process request body",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}

	if err := json.Unmarshal(body, &entropy); err != nil {
		var s = "invalid request object, expected json format"
		handleError( w, r, http.StatusUnprocessableEntity, s )
		return
	}

	if ok, err := entropy.Validate(); !ok {
		logger.Error( "validating request object",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	var score response.Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	score, err = calc.Calc( reqCtx, &entropy, formatId, useReductions )
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		logger.Error( "calculating request score",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

	} else {
		ds, err := scoringdb.GetDataStore( nil )
		if err == nil {
			reqId, _ := logging.GetReqId( reqCtx )
			rec, _ := scoringdb.NewReqRecord( &entropy, reqId, time.Now() )

			err := ds.SaveRequest( reqCtx, rec )
			if err != nil {
				logger.Error( "logging request score",
					zap.String("formatId", strings.ToLower(formatId) ),
					zap.String( "status", "error" ),
					zap.String("error ", err.Error() ),
				)
			}

			resp, _ := scoringdb.NewRespRecord( &score, reqId, time.Now() )

			err = ds.SaveResponse( reqCtx, resp )
			if err != nil {
				logger.Error( "logging response score",
					zap.String("formatId", strings.ToLower(formatId) ),
					zap.String( "status", "error" ),
					zap.String("error ", err.Error() ),
				)
			}
		}

		w.WriteHeader(http.StatusOK)
		if modeId == "summary" {
			score.Data.People = nil
		}
	}

	if err := json.NewEncoder(w).Encode(score); err != nil {
		logger.Error( "encoding scoring response",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	var score response.Response

	reqCtx := logging.WithFuncId( r.Context(), "handleError", "scores" )
	logger := logging.Logger(reqCtx)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader( statusCode )

	score.Errors = new(response.Errors)
	score.Errors.Messages = append( score.Errors.Messages, msg )
	if err := json.NewEncoder(w).Encode(score); err != nil {
		logger.Error( "handling scoring error",
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}
}

