package scores

import (
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"context"
	"regexp"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/request"
	"github.com/mfioravanti2/entropy-api/model/response"
	"github.com/mfioravanti2/entropy-api/calc"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

const (
	DEFAULT_SCORING = "mean"
)

func Validate( formatId string ) (bool, error) {
	var err error

	rx, err := regexp.Compile( `^(mean|naive|rare)$` )
	if err != nil {
		return false, err
	}

	if rx.MatchString( strings.ToLower(formatId) ) {
		return true, nil
	}

	return false, nil
}

func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "scores" )

	logger := logging.Logger( ctx )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores" ) )
	r = append( r, model.Route{"ScoreCalc", "POST", "/v1/scores", Calc} )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/scores/format/{formatId}" ) )
	r = append( r, model.Route{"ScoreCalcFormat", "POST", "/v1/scores/format/{formatId}", CalcFormat} )

	return r
}

func Calc(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	logger.Info( "preparing to score request, with default formatId",
		zap.String("formatId", strings.ToLower(DEFAULT_SCORING) ),
	)

	score( w, r, DEFAULT_SCORING)
}

func CalcFormat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	formatId := strings.ToLower(vars["formatId"])

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	logger.Info( "preparing to score request, with request specified formatId",
		zap.String("formatId", strings.ToLower(formatId) ),
	)

	score( w, r, formatId )
}

func score(w http.ResponseWriter, r *http.Request, formatId string) {
	var entropy request.Request

	reqCtx := logging.WithFuncId( r.Context(), "score", "scores" )
	logger := logging.Logger(reqCtx)

	if ok, _ := Validate(formatId); !ok {
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

	var score response.Response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if score, err = calc.Calc( reqCtx, &entropy, formatId ); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		logger.Error( "calculating request score",
			zap.String("formatId", strings.ToLower(formatId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

	} else {
		w.WriteHeader(http.StatusOK)
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

