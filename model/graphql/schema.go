package entropyql

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

var entropySchema *graphql.Schema

func init() {
	BuildSchema()
}

func BuildSchema() {
	ctx := logging.WithFuncId( context.Background(), "BuildSchema", "entropyql" )

	logger := logging.Logger( ctx )
	logger.Debug("building GraphQL schema", zap.String( "type", "schema" ) )

	var countryType *graphql.Object
	countryType = getCountryType()

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"countries": &graphql.Field{
				Type: graphql.NewList(countryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.countries" )
					ctrReg.Inc(1)

					return data.GetAllCountries(), nil
				},
			},
		},
	})

	var schema graphql.Schema
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
	})

	entropySchema = &schema
}

func GetSchema() ( *graphql.Schema, error ) {
	if entropySchema != nil {
		return entropySchema, nil
	}

	s := fmt.Sprintf( "unable to load schema, nil schema" )
	return nil, errors.New( s )
}
