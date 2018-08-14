package entropyql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model/metrics"
)

var entropySchema *graphql.Schema

func init() {
	BuildSchema()
}

func BuildSchema() {
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
