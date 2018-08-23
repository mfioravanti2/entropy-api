package entropyql

import (
	"context"

	"go.uber.org/zap"
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func getHeuristicType() *graphql.Object {
	ctx := logging.WithFuncId( context.Background(), "getHeuristicType", "entropyql" )

	logger := logging.Logger( ctx )
	logger.Debug("building GraphQL schema",
		zap.String( "type", "heuristicType" ),
		)

	var heuristicType *graphql.Object

	heuristicType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Heuristic",
		Description: "A heuristic for a country model.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Globally unique identifier for the Heuristic (UUID v4 format)",
				Resolve: resolveHeuristicId,
			},
			"notes": &graphql.Field{
				Type: graphql.String,
				Description: "Description of the Heuristic",
				Resolve: resolveHeuristicNotes,
			},
			"match": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Description: "Mnemonics which must be matched before this heuristic is applied",
				Resolve: resolveHeuristicMatch,
			},
			"insert": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Description: "Mnemonics that are removed when this heuristic is triggered",
				Resolve: resolveHeuristicInsert,
			},
			"remove": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Description: "Mnemonics that are removed when this heuristic is triggered",
				Resolve: resolveHeuristicRemove,
			},
		},
	})

	return heuristicType
}

