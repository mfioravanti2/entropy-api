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

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/attributes" ) )
	r = append( r, model.Route{"AttributeList", "GET", "/v1/countries/{countryId}/attributes", List, nil})

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/attributes/{attributeId}" ) )
	r = append( r, model.Route{"AttributeDetails", "GET", "/v1/countries/{countryId}/attributes/{attributeId}", Detail, nil})

	return r
}

func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	if ok, _ := model.ValidateCountryCode(countryId); !ok {
		logger.Error( "validating country code",
			zap.String("countryId", strings.ToUpper(countryId)),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid country code" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	logger.Info( "retrieving attributes from country model",
		zap.String("countryId", strings.ToUpper(countryId) ),
	)

	var err error
	var attributes []string
	attributes, err = data.GetAttributes(countryId)
	if err != nil {
		logger.Error( "retrieving attributes from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if len(attributes) > 0 {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(attributes); err != nil {
			logger.Error( "encoding attributes",
				zap.String("countryId", strings.ToUpper(countryId) ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )

		logger.Info( "retrieving attributes from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "status", "ok" ),
			zap.String("error ", "no attributes found" ),
		)
	}
}

func Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])
	attributeId := strings.ToLower(vars["attributeId"])

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	if ok, _ := model.ValidateCountryCode(countryId); !ok {
		logger.Error( "validating country code",
			zap.String("countryId", strings.ToUpper(countryId)),
			zap.String( "attributeId", strings.ToLower(attributeId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid country code" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	if ok, _ := model.ValidateAttributeMnemonic(attributeId); !ok {
		logger.Error( "validating attribute identifier",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "attributeId", strings.ToLower(attributeId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid attribute identifier" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )
		return
	}

	logger.Debug( "retrieving attribute from country model",
		zap.String("countryId", strings.ToUpper(countryId) ),
		zap.String( "attributeId", strings.ToLower(attributeId) ),
	)

	attribute, err := data.GetAttribute(countryId, attributeId)

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if err == nil {
		w.WriteHeader( http.StatusOK )

		if err := json.NewEncoder(w).Encode(attribute); err != nil {
			logger.Error( "encoding attribute",
				zap.String("countryId", strings.ToUpper(countryId) ),
				zap.String( "attributeId", strings.ToLower(attributeId) ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}
	} else {
		w.WriteHeader( http.StatusNotFound )

		logger.Error( "retrieving attribute from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "attributeId", strings.ToLower(attributeId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}
}
