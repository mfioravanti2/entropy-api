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
	"github.com/mfioravanti2/entropy-api/cli"
)

// Add Handlers for the Country Endpoints
func AddHandlers(r model.Routes, endpoints *cli.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "country" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( cli.ENDPOINT_REST )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", cli.ENDPOINT_REST ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries" ) )
			r = append( r, model.Route{"CountryList", "GET", "/v1/countries", List, nil} )
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/countries" ),
				zap.String( "policy", cli.ENDPOINT_REST ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", cli.ENDPOINT_REST ),
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
		w.Header().Set("Content-type", "application/json; charset=UTF-8")
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

