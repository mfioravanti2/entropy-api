package graphql

import (
	"context"
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"

	"go.uber.org/zap"
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/model/graphql"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/command/server/enforce"
)

// Add Handlers for the GraphQL Endpoints
func AddHandlers(r model.Routes, endpoints *config.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "graphql" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( config.ENDPOINT_GRAPHQL )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", config.ENDPOINT_GRAPHQL ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers",
				zap.String( "endpoint", "/v1/graphql" ),
				)
			r = append( r, model.Route{ Name: "GraphQL",
										Method: "POST",
										Pattern: "/v1/graphql",
										HandlerFunc: GraphQL,
										Params: nil,
										Enforce: model.ENFORCE_CONTENT_GRAPHQL,
										Policy: endpoint,
										AuthN: model.AUTH_METHOD_NONE  } )
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/graphql" ),
				zap.String( "policy", config.ENDPOINT_GRAPHQL ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", config.ENDPOINT_GRAPHQL ),
		)
	}

	return r
}

// Return the GraphQL object
func GraphQL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", enforce.HEADER_JSON_CONTENT_TYPE)

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.graphql.post" )
	ctrReg.Inc(1)

	logger.Debug("preparing to parse query graphql message",
	)

	// Read the body from the request body
	// Maximum request by size is 50kb
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 50 * 1024 ))
	if err != nil {
		w.WriteHeader( http.StatusUnprocessableEntity )

		logger.Error( "unable to read request body",
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.422" )
		ctrReg.Inc(1)

		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader( http.StatusUnprocessableEntity )

		logger.Error( "unable to process request body",
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.422" )
		ctrReg.Inc(1)

		return
	}

	logger.Info( "parsing graphql query/mutation",
	)

	var gqlQuery map[string]interface{}
	if err := json.Unmarshal(body, &gqlQuery); err != nil {
		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.422" )
		ctrReg.Inc(1)

		return
	}

	// query is type string by default
	var query string
	// variables is nil if no variables are supplied, or of type
	// map[string]interface {} if variables have been supplied
	var variables interface{}

	var entropySchema *graphql.Schema

	var result interface{}

	if entropySchema, err = entropyql.GetSchema(); err != nil {
		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.422" )
		ctrReg.Inc(1)

		return
	}

	query = gqlQuery["query"].(string)
	variables = gqlQuery["variables"]

	if variables == nil {
		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.variables" )
		ctrReg.Inc(1)

		// No variables have been provided, call the GraphQL library
		// without assigning a variable parameter
		logger.Info( "evaluating graphql query/mutation",
			zap.Bool( "params", false ),
		)

		result = graphql.Do(graphql.Params{
			Schema:         *entropySchema,
			RequestString:  query,
		})
	} else {
		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.no_variables" )
		ctrReg.Inc(1)

		// Variables have been provided, call the GraphQL library
		// with the specified variable parameter
		logger.Info( "evaluating graphql query/mutation",
			zap.Bool( "params", true ),
		)

		result = graphql.Do(graphql.Params{
			Schema:         *entropySchema,
			RequestString:  query,
			VariableValues: variables.(map[string]interface{}),
		})
	}

	// Encode and return the graphql response
	w.Header().Set("Content-Type", enforce.HEADER_JSON_CONTENT_TYPE)
	w.WriteHeader( http.StatusOK )
	if err := json.NewEncoder(w).Encode(result); err != nil {
		logger.Error( "encoding graphql response",
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}

	if err == nil {
		logger.Info( "returning graphql object",
			zap.String( "status", "ok" ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.200" )
		ctrReg.Inc(1)
	} else {
		logger.Error( "returning graphql object",
			zap.String( "status", "error" ),
			zap.String( "error", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.500" )
		ctrReg.Inc(1)
	}
}
