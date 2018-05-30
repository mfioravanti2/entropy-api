package attribute

import (
	"net/http"
	"encoding/json"
	"strings"

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
	countryId := strings.ToLower(vars["countryId"])

	var attributes []string
	attributes = data.GetAttributes(countryId)

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if len(attributes) > 0 {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(attributes); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}
