package country

import (
	"encoding/json"
	"net/http"
	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
)

func AddHandlers(r model.Routes) model.Routes {
	r = append( r, model.Route{"CountryList", "GET", "/v1/countries", List} )

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader( http.StatusOK )

	if err := json.NewEncoder(w).Encode(data.GetCountries()); err != nil {
		panic(err)
	}
}

