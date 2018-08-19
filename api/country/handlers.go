package country

import (
	"context"
	"net/http"
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"

	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/command/server/enforce"
)

// Add Handlers for the Country Endpoints
func AddHandlers(r model.Routes, endpoints *config.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "country" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( config.ENDPOINT_REST )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", config.ENDPOINT_REST ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers",
				zap.String( "endpoint", "/v1/countries" ),
				)
			r = append( r, model.Route{ Name: "CountryList",
										Method: "GET",
										Pattern: "/v1/countries",
										HandlerFunc: List,
										Params: nil,
										Enforce: model.ENFORCE_CONTENT_NONE,
										Policy: endpoint,
										AuthN: model.AUTH_METHOD_NONE } )
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/countries" ),
				zap.String( "policy", config.ENDPOINT_REST ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", config.ENDPOINT_REST ),
		)
	}

	return r
}

// List the country codes available
func List(w http.ResponseWriter, r *http.Request) {
	var countries []string
	countries = data.GetCountries()

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.country.get" )
	ctrReg.Inc(1)

	logger.Info( "retrieving country codes from models" )

	if len(countries) > 0 {
		w.Header().Set("Content-Type", enforce.HEADER_JSON_CONTENT_TYPE)
		w.WriteHeader( http.StatusOK )

		// obtain a list of country codes
		for i, countryCode := range countries {
			countries[i] = strings.ToUpper( countryCode )
		}

		// encode and return the array of country codes
		if err := json.NewEncoder(w).Encode(countries); err != nil {
			logger.Error( "encoding country codes",
				zap.String( "status", "error" ),
				zap.String("error", err.Error() ),
			)
		} else {
			logger.Info( "retrieved country codes from models",
				zap.String( "status", "ok" ),
			)
		}

		ctrReg, _ := metrix.GetCounter( "entropy.country.get.status.200" )
		ctrReg.Inc(1)
	} else {
		w.WriteHeader( http.StatusNoContent )

		logger.Info( "retrieved country codes from models",
			zap.String( "status", "ok" ),
			zap.String("error ", "no country codes found" ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.country.get.status.204" )
		ctrReg.Inc(1)
	}
}

