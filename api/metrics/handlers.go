package metrics

import (
	"context"
	"net/http"
	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"

	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/command/server/enforce"
)

// Add Handlers for the Metrics-based Endpoints
func AddHandlers(r model.Routes, endpoints *config.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "metrics" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( config.ENDPOINT_METRICS )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", config.ENDPOINT_METRICS ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers", zap.String( "endpoint", "/v1/sys/metrics" ) )
			r = append( r, model.Route{ Name: "SysMetrics", Method: "GET", Pattern: "/v1/sys/metrics", HandlerFunc: Metrics } )
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/sys/metrics" ),
				zap.String( "policy", config.ENDPOINT_METRICS ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", config.ENDPOINT_METRICS ),
		)
	}


	return r
}

// Return the Metrics as a json object
func Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", enforce.HEADER_JSON_CONTENT_TYPE)

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.sys.metrics.get" )
	ctrReg.Inc(1)

	logger.Debug("preparing to retrieve metrics",
	)

	jsonMetrics, err := metrix.GetJson()

	if err == nil {
		w.WriteHeader( http.StatusOK )
		w.Write(jsonMetrics)

		logger.Info( "returning metrics object",
			zap.String( "status", "ok" ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.sys.metrics.get.status.200" )
		ctrReg.Inc(1)
	} else {
		w.WriteHeader( http.StatusInternalServerError )

		logger.Error( "returning metrics object",
			zap.String( "status", "error" ),
			zap.String( "error", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.sys.metrics.get.status.500" )
		ctrReg.Inc(1)
	}
}
