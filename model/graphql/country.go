package entropyql

import (
	"context"

	"go.uber.org/zap"
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func getCountryType() *graphql.Object {
	ctx := logging.WithFuncId( context.Background(), "getCountryType", "entropyql" )

	logger := logging.Logger( ctx )
	logger.Debug("building GraphQL schema",
		zap.String( "type", "countryType" ),
		)

	var countryType *graphql.Object

	countryType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Country",
		Description: "A set of country models.",
		Fields: graphql.Fields{
			"locale": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "A 2-digit Country Code as defined in ISO 3166-1 alpha-2",
				Resolve: resolveCountryLocale,
			},
			"threshold": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
				Description: "Entropy Threshold for the Model",
				Resolve: resolveCountryThreshold,
			},
			"k": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Description: "k-anonymity value for the Model",
				Resolve: resolveCountryK,
			},
			"version": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Country Model Version",
				Resolve: resolveCountryVersion,
			},
			"timestamp": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
				Description: "Country Code for the Model",
				Resolve: resolveCountryTimestamp,
			},
			"heuristics": &graphql.Field{
				Type: graphql.NewList( getHeuristicType() ),
				Description: "List of Heuristics which operate on the Country Model's Attributes",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
						Description: "Globally unique identifier for the Heuristic (UUID v4 format)",
					},
				},
				Resolve: resolveCountryHeuristics,
			},
			"attributes": &graphql.Field{
				Type: graphql.NewNonNull( graphql.NewList(getAttributeType() )),
				Description: "List of Scored Attributes associated with the Model",
				Args: graphql.FieldConfigArgument{
					"mnemonic": &graphql.ArgumentConfig{
						Type: graphql.String,
						Description: "Country Model unique identifier for the Attribute",
					},
				},
				Resolve: resolveCountryAttributes,
			},
		},
	})

	return countryType
}
