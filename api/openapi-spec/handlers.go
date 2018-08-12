package openapi_spec

import (
	"context"
	"net/http"
	"io/ioutil"
	"fmt"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

// Add Handlers for the Specification-based Endpoints
func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "sys" )

	logger := logging.Logger( ctx )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/sys/spec" ) )
	r = append( r, model.Route{ Name: "SysSchema", Method: "GET", Pattern: "/v1/sys/spec", HandlerFunc: Spec} )

	return r
}

// Return the OpenAPI specification as a json object
func Spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	logger.Debug("preparing to reloading models",
		zap.String("model_file", "data/sources/sources.json" ),
	)

	jsonSpec, err := ioutil.ReadFile("api/openapi-spec/openapi.json")
	if err != nil {
		s := fmt.Sprintf("unable to load OpenAPI specification file")
		logger.Error( "loading openapi spec",
			zap.String("file", "api/openapi-spec/openapi.json" ),
			zap.String("error", s ),
		)
	}

	if err == nil {
		w.WriteHeader( http.StatusOK )
		w.Write(jsonSpec)

		logger.Info( "returning openapi spec",
			zap.String( "status", "ok" ),
		)
	} else {
		w.WriteHeader( http.StatusInternalServerError )

		logger.Error( "returning openapi spec",
			zap.String( "status", "error" ),
			zap.String( "error", err.Error() ),
		)
	}
}
