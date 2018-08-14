package entropyql

import (
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model/source"
)

func getCountryType() *graphql.Object {
	var countryType *graphql.Object

	countryType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Country",
		Description: "A set of country models.",
		Fields: graphql.Fields{
			"locale": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Country Code for the Model",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Model); ok {
						return model.Locale, nil
					}
					return nil, nil
				},
			},
			"threshold": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
				Description: "Entropy Threshold for the Model",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Model); ok {
						return model.Threshold, nil
					}
					return nil, nil
				},
			},
			"k": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
				Description: "k-anonymity value for the Model",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Model); ok {
						return model.K, nil
					}
					return nil, nil
				},
			},
			"version": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Country Model Version",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Model); ok {
						return model.ModelVersion, nil
					}
					return nil, nil
				},
			},
			"timestamp": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
				Description: "Country Code for the Model",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Model); ok {
						return model.ModelDate, nil
					}
					return nil, nil
				},
			},
			"heuristics": &graphql.Field{
				Type: graphql.NewList( getHeuristicType() ),
				Description: "List of Heuristics which operate on the Country Model's Attributes",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Model); ok {
						return model.Heuristics, nil
					}
					return nil, nil
				},
			},
			"attributes": &graphql.Field{
				Type: graphql.NewNonNull( graphql.NewList(getAttributeType() )),
				Description: "List of Scored Attributes associated with the Model",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if model, ok := p.Source.(source.Model); ok {
						return model.Attributes, nil
					}
					return nil, nil
				},
			},
		},
	})

	return countryType
}
