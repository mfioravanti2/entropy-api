package attribute

import (
	"context"
	"net/http"
	"encoding/json"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "sys" )

	logger := logging.Logger( ctx )

	logger.Info("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/attributes" ) )
	r = append( r, model.Route{"AttributeList", "GET", "/v1/countries/{countryId}/attributes", List})

	logger.Info("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/attributes/{attributeId}" ) )
	r = append( r, model.Route{"AttributeDetails", "GET", "/v1/countries/{countryId}/attributes/{attributeId}", Detail})

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])

	var err error
	var attributes []string
	attributes, err = data.GetAttributes(countryId)
	if err != nil {
		panic(err)
	}

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

func Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])
	attributeId := strings.ToLower(vars["attributeId"])

	attribute, err := data.GetAttribute(countryId, attributeId)

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if err == nil {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(attribute); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )
	}
}
