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
	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/cli"
)

// Add Handlers for the Attribute Endpoints
func AddHandlers(r model.Routes, endpoints *cli.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "attribute" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( cli.ENDPOINT_REST )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", cli.ENDPOINT_REST ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/attributes" ) )
			r = append( r, model.Route{"AttributeList", "GET", "/v1/countries/{countryId}/attributes", List, nil})

			logger.Debug("registering handlers", zap.String( "endpoint", "/v1/countries/{countryId}/attributes/{attributeId}" ) )
			r = append( r, model.Route{"AttributeDetails", "GET", "/v1/countries/{countryId}/attributes/{attributeId}", Detail, nil})
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/countries/{countryId}/attributes" ),
				zap.String( "policy", cli.ENDPOINT_REST ),
				)

			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/countries/{countryId}/attributes/{attributeId}" ),
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

// List the Attributes associated with a specified country
func List(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.attributes_list.get" )
	ctrReg.Inc(1)

	// validate the country code
	if ok, _ := model.ValidateCountryCode(countryId); !ok {
		logger.Error( "validating country code",
			zap.String("countryId", strings.ToUpper(countryId)),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid country code" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.attributes_list.get.status.422" )
		ctrReg.Inc(1)
		return
	}

	logger.Info( "retrieving attributes from country model",
		zap.String("countryId", strings.ToUpper(countryId) ),
	)

	var err error
	var attributes []string

	// retrieve a list of attributes from the specified country's model
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

		// encode and return a list of attributes available within the country's model
		if err := json.NewEncoder(w).Encode(attributes); err != nil {
			logger.Error( "encoding attributes",
				zap.String("countryId", strings.ToUpper(countryId) ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}

		ctrReg, _ := metrix.GetCounter( "entropy.attributes_list.get.status.200" )
		ctrReg.Inc(1)
	} else {
		w.WriteHeader( http.StatusNotFound )

		logger.Info( "retrieving attributes from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "status", "ok" ),
			zap.String("error ", "no attributes found" ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.attributes_list.get.status.404" )
		ctrReg.Inc(1)
	}
}

// Provide details about a specific attribute from a country's model
func Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryId := strings.ToLower(vars["countryId"])
	attributeId := strings.ToLower(vars["attributeId"])

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.attributes_details.get" )
	ctrReg.Inc(1)

	// Validate the country code
	if ok, _ := model.ValidateCountryCode(countryId); !ok {
		logger.Error( "validating country code",
			zap.String("countryId", strings.ToUpper(countryId)),
			zap.String( "attributeId", strings.ToLower(attributeId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid country code" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.attributes_details.get.status.422" )
		ctrReg.Inc(1)
		return
	}

	// Validate the attribute's mnemonic format
	if ok, _ := model.ValidateAttributeMnemonic(attributeId); !ok {
		logger.Error( "validating attribute identifier",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "attributeId", strings.ToLower(attributeId) ),
			zap.String( "status", "error" ),
			zap.String("error ", "invalid attribute identifier" ),
		)

		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.attributes_details.get.status.422" )
		ctrReg.Inc(1)
		return
	}

	logger.Debug( "retrieving attribute from country model",
		zap.String("countryId", strings.ToUpper(countryId) ),
		zap.String( "attributeId", strings.ToLower(attributeId) ),
	)

	// Get information about the specified attribute from the country's model
	attribute, err := data.GetAttribute(countryId, attributeId)

	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	if err == nil {
		w.WriteHeader( http.StatusOK )

		// Encode and return the attribute
		if err := json.NewEncoder(w).Encode(attribute); err != nil {
			logger.Error( "encoding attribute",
				zap.String("countryId", strings.ToUpper(countryId) ),
				zap.String( "attributeId", strings.ToLower(attributeId) ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}

		ctrReg, _ := metrix.GetCounter( "entropy.attributes_details.get.status.200" )
		ctrReg.Inc(1)
	} else {
		w.WriteHeader( http.StatusNotFound )

		logger.Error( "retrieving attribute from country model",
			zap.String("countryId", strings.ToUpper(countryId) ),
			zap.String( "attributeId", strings.ToLower(attributeId) ),
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.attributes_details.get.status.404" )
		ctrReg.Inc(1)
	}
}
