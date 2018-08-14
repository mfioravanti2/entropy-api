package entropyql

import (
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model/source"
)

func getHeuristicType() *graphql.Object {
	var heuristicType *graphql.Object

	heuristicType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Heuristic",
		Description: "A heuristic for a country model.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Globally unique identifier for the Heuristic (UUID v4 format)",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Heuristic); ok {
						return model.Id, nil
					}
					return nil, nil
				},
			},
			"notes": &graphql.Field{
				Type: graphql.String,
				Description: "Description of the Heuristic",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Heuristic); ok {
						return model.Notes, nil
					}
					return nil, nil
				},
			},
			"match": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Description: "Mnemonics which must be matched before this heuristic is applied",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Heuristic); ok {
						return model.Match, nil
					}
					return nil, nil
				},
			},
			"insert": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Description: "Mnemonics that are removed when this heuristic is triggered",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Heuristic); ok {
						return model.Insert, nil
					}
					return nil, nil
				},
			},
			"remove": &graphql.Field{
				Type: graphql.NewList(graphql.String),
				Description: "Mnemonics that are removed when this heuristic is triggered",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Heuristic); ok {
						return model.Remove, nil
					}
					return nil, nil
				},
			},
		},
	})

	return heuristicType
}

