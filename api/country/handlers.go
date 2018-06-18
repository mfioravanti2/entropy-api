package country

import (
	"net/http"
	"encoding/json"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func AddHandlers(r model.Routes) model.Routes {
	logger := logging.Logger(nil)
	logger.Info("AddHandlers(country)")

	r = append( r, model.Route{"CountryList", "GET", "/v1/countries", List} )

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	var countries []string
	countries = data.GetCountries()

	if len(countries) > 0 {
		w.Header().Set("Content-type", "application/json; charset=UTF-8")
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(countries); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader( http.StatusNoContent )
	}
}

