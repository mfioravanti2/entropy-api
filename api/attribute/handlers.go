package attribute

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
)

func AddHandlers(r model.Routes) model.Routes {
	r = append( r, model.Route{"AttributeList", "GET", "/v1/countries/{countryId}/attributes", List})

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := vars["countryId"]

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader( http.StatusOK )

	if err := json.NewEncoder(w).Encode(data.GetAttributes(countryId)); err != nil {
		panic(err)
	}
}
