package graphql

import (
	"context"
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"

	"go.uber.org/zap"
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/command/server/logging"

	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/model/graphql"
)

// Add Handlers for the GraphQL Endpoints
func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "graphql" )

	logger := logging.Logger( ctx )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/graphql" ) )
	r = append( r, model.Route{ Name: "GraphQL", Method: "POST", Pattern: "/v1/graphql", HandlerFunc: GraphQL } )

	return r
}

// Return the GraphQL object
func GraphQL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.graphql.post" )
	ctrReg.Inc(1)

	logger.Debug("preparing to query graphql",
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

	var gqlQuery map[string]interface{}
	if err := json.Unmarshal(body, &gqlQuery); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.422" )
		ctrReg.Inc(1)

		return
	}

	query := gqlQuery["query"]
//	variables := gqlQuery["variables"]
	var entropySchema *graphql.Schema

	if entropySchema, err = entropyql.GetSchema(); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader( http.StatusUnprocessableEntity )

		ctrReg, _ := metrix.GetCounter( "entropy.graphql.post.status.422" )
		ctrReg.Inc(1)

		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         *entropySchema,
		RequestString:  query.(string),
//		VariableValues: variables.(map[string]interface{}),
	})

	// Encode and return the graphql response
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
