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
	"github.com/mfioravanti2/entropy-api/model/metrics"
	"fmt"
	"github.com/mfioravanti2/entropy-api/cli"
)

const (
	// valid values are naive, mean, or rare
	DEFAULT_SCORING = "mean"
	// valid values are detailed or summary
	DEFAULT_MODE = "detailed"
	// valide values are exclude or include
	DEFAULT_REDUCTIONS = "include"
)

//	Generate a complete list of available routes
func AddHandlers(r model.Routes, endpoints *cli.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "scores" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( cli.ENDPOINT_SCORING )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", cli.ENDPOINT_SCORING ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
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
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/scores" ),
				zap.String( "policy", cli.ENDPOINT_SCORING ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", cli.ENDPOINT_SCORING ),
		)
	}

	return r
}

// Process an attribute set with the default options
// 	formatId: mean, use the arithmetic mean for calculating the attribute's entropy scores
// 	modeId: detailed, return details results once the scoring process has completed in the response
//  reductions: true, apply reduction heuristics during the scoring process
func CalcDefaults(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.scoring.post" )
	ctrReg.Inc(1)

	ctrReg, _ = metrix.GetCounter( "entropy.scoring.default" )
	ctrReg.Inc(1)

	logger.Info( "preparing to score request, with default formatId",
		zap.String("formatId", strings.ToLower(DEFAULT_SCORING) ),
	)

	// collect scoring options and process the request
	score( w, r, DEFAULT_MODE, DEFAULT_SCORING, true )
}

// Process an attribute set with client specified scoring options
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

	ctrReg, _ := metrix.GetCounter( "entropy.scoring.post" )
	ctrReg.Inc(1)

	ctrReg, _ = metrix.GetCounter( "entropy.scoring.options" )
	ctrReg.Inc(1)

	logger.Info( "preparing to score request, with request specified formatId",
		zap.String("formatId", strings.ToLower(formatId) ),
		zap.String("modeId", strings.ToLower(modeId) ),
		zap.Bool("useReductions", useReductions == "include" ),
	)

	// collect scoring options and process the request
	score( w, r, modeId, formatId, useReductions == "include" )
}

// Score an attribute set
func score(w http.ResponseWriter, r *http.Request, modeId string, formatId string, useReductions bool ) {
	var entropy request.Request

	reqCtx := logging.WithFuncId( r.Context(), "score", "scores" )
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.scoring.mode." + modeId )
	ctrReg.Inc(1)

	// Validate the formatId
	if ok, _ := model.ValidateFormat(formatId); !ok {
		logger.Error( "validating format identifier",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid format identifier" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.scoring.post.status.422" )
		ctrReg.Inc(1)
		return
	}

	logger.Info( "score request",
		zap.String("formatId", strings.ToLower(formatId) ),
		zap.String("modeId", strings.ToLower(modeId) ),
		zap.Bool("useReductions", useReductions ),
	)

	// Read the body from the request body
	// Maximum request by size is 50kb
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

	// Convert the request body into a scoring object
	if err := json.Unmarshal(body, &entropy); err != nil {
		var s = "invalid request object, expected json format"
		handleError( w, r, http.StatusUnprocessableEntity, s )
		return
	}

	// Validate the scoring request
	if ok, err := entropy.Validate(); !ok {
		logger.Error( "validating request object",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.scoring.post.status.422" )
		ctrReg.Inc(1)
		return
	}

	var score response.Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Score the request
	score, err = calc.Calc( reqCtx, &entropy, formatId, useReductions )
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		logger.Error( "calculating request score",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.scoring.post.status.422" )
		ctrReg.Inc(1)
	} else {
		// Retrieve the Request/Response Data Store
		ds, err := scoringdb.GetDataStore( nil )
		if err == nil {
			reqId, _ := logging.GetReqId( reqCtx )
			// Convert the scoring request
			rec, _ := scoringdb.NewReqRecord( &entropy, reqId, time.Now() )

			// Record the Scoring Request
			err := ds.SaveRequest( reqCtx, rec )
			if err != nil {
				logger.Error( "logging request score",
					zap.String("formatId", strings.ToLower(formatId) ),
					zap.String( "status", "error" ),
					zap.String("error ", err.Error() ),
				)
			}

			// Convert the scoring response
			resp, _ := scoringdb.NewRespRecord( &score, reqId, time.Now() )

			// Recording the Scoring Response
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
			// If the client is expecting only the scoring summary, remove the details from the response
			score.Data.People = nil
		}

		ctrReg, _ := metrix.GetCounter( "entropy.scoring.post.status.200" )
		ctrReg.Inc(1)
	}

	// Encode and return the scoring response
	if err := json.NewEncoder(w).Encode(score); err != nil {
		logger.Error( "encoding scoring response",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}
}

// Handle an Error and return a custom error response
func handleError(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	var score response.Response

	reqCtx := logging.WithFuncId( r.Context(), "handleError", "scores" )
	logger := logging.Logger(reqCtx)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader( statusCode )

	cLabel := fmt.Sprintf( "entropy.scoring.post.status.%d", statusCode )
	ctrReg, _ := metrix.GetCounter( cLabel )
	ctrReg.Inc(1)

	score.Errors = new(response.Errors)
	score.Errors.Messages = append( score.Errors.Messages, msg )
	if err := json.NewEncoder(w).Encode(score); err != nil {
		logger.Error( "handling scoring error",
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}
}

