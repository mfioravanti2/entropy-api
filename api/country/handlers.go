package country

import (
	"context"
	"net/http"
	"encoding/json"
	"regexp"
	"strings"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "country" )

	logger := logging.Logger( ctx )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries" ) )
	r = append( r, model.Route{"CountryList", "GET", "/v1/countries", List, nil} )

	return r
}

func Validate( countryCode string ) (bool, error) {
	var err error

	// two-digit country codes: ISO 3166-1 alpha-2
	rx, err := regexp.Compile(`^[a-zA-Z]{2}$` )
	if err != nil {
		return false, err
	}

	if rx.MatchString( countryCode ) {
		return true, nil
	}

	return false, nil
}

func List(w http.ResponseWriter, r *http.Request) {
	var countries []string
	countries = data.GetCountries()

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	logger.Info( "retrieving country codes from models" )

	if len(countries) > 0 {
		w.Header().Set("Content-type", "application/json; charset=UTF-8")
		w.WriteHeader( http.StatusOK )

		for i, countryCode := range countries {
			countries[i] = strings.ToUpper( countryCode )
		}

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
	} else {
		w.WriteHeader( http.StatusNoContent )

		logger.Info( "retrieved country codes from models",
			zap.String( "status", "ok" ),
			zap.String("error ", "no country codes found" ),
		)
	}
}

