package openapi_spec

import (
	"context"
	"net/http"
	"io/ioutil"
	"fmt"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/cli"
)

// Add Handlers for the Specification-based Endpoints
func AddHandlers(r model.Routes, endpoints *cli.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "openapi_spec" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( cli.ENDPOINT_OPENAPI )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", cli.ENDPOINT_OPENAPI ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers", zap.String( "endpoint", "/v1/sys/spec" ) )
			r = append( r, model.Route{ Name: "SysSchema", Method: "GET", Pattern: "/v1/sys/spec", HandlerFunc: Spec} )
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/sys/spec" ),
				zap.String( "policy", cli.ENDPOINT_OPENAPI ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", cli.ENDPOINT_OPENAPI ),
		)
	}

	return r
}

// Return the OpenAPI specification as a json object
func Spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.sys.openapi_spec.get" )
	ctrReg.Inc(1)

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

		ctrReg, _ := metrix.GetCounter( "entropy.sys.openapi_spec.get.status.200" )
		ctrReg.Inc(1)
	} else {
		w.WriteHeader( http.StatusInternalServerError )

		logger.Error( "returning openapi spec",
			zap.String( "status", "error" ),
			zap.String( "error", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.sys.openapi_spec.get.status.404" )
		ctrReg.Inc(1)
	}
}
