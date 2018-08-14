package entropyql

import (
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model/source"
	"github.com/mfioravanti2/entropy-api/model/metrics"
)

func getAttributeFormatType() *graphql.Object {
	var formatType *graphql.Object

	formatType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Format",
		Description: "Calculation Format of a Scored Attribute",
		Fields: graphql.Fields{
			"format": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Method of Calculation for the Scored Value",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.format.method" )
					ctrReg.Inc(1)

					if f, ok := p.Source.(source.Format); ok {
						return f.Format, nil
					}

					return nil, nil
				},
			},
			"score": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
				Description: "Scored Value of an Attribute",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.format.score" )
					ctrReg.Inc(1)

					if f, ok := p.Source.(source.Format); ok {
						return f.Score, nil
					}

					return nil, nil
				},
			},
		},
	})

	return formatType
}

func getAttributeSourceType() *graphql.Object {
	var sourceType *graphql.Object

	sourceType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Source",
		Description: "Source of Frequency Statistics used to perform a Scoring Calculation",
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Document Title",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if s, ok := p.Source.(source.Source); ok {
						return s.Title, nil
					}

					return nil, nil
				},
			},
			"date": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
				Description: "Date of Publication or Retrieval",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if s, ok := p.Source.(source.Source); ok {
						return s.Date, nil
					}

					return nil, nil
				},
			},
			"organization": &graphql.Field{
				Type: graphql.String,
				Description: "Organization Responsible for the Source Document",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if s, ok := p.Source.(source.Source); ok {
						return s.Org, nil
					}

					return nil, nil
				},
			},
			"url": &graphql.Field{
				Type: graphql.String,
				Description: "URL for retrieving the Source Document",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if s, ok := p.Source.(source.Source); ok {
						return s.URI, nil
					}

					return nil, nil
				},
			},
		},
	})

	return sourceType
}

func getAttributeType() *graphql.Object {
	var attributeType *graphql.Object

	attributeType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Attribute",
		Description: "An attribute for a Country Model.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Globally unique identifier for the Attribute (UUID v4 format)",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.id" )
					ctrReg.Inc(1)

					if a, ok := p.Source.(source.Attribute); ok {
						return a.Id, nil
					}

					return nil, nil
				},
			},
			"mnemonic": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Country Model unique identifier for the Attribute",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.mnemonic" )
					ctrReg.Inc(1)

					if a, ok := p.Source.(source.Attribute); ok {
						return a.Mnemonic, nil
					}

					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Description: "Human friendly short description of the Attribute",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if a, ok := p.Source.(source.Attribute); ok {
						return a.Name, nil
					}

					return nil, nil
				},
			},
			"notes": &graphql.Field{
				Type: graphql.String,
				Description: "Description of the Attribute",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if a, ok := p.Source.(source.Attribute); ok {
						return a.Notes, nil
					}

					return nil, nil
				},
			},
			"sources": &graphql.Field{
				Type: graphql.NewList( getAttributeSourceType() ),
				Description: "Source Documents for the Scoring Calculations",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.sources" )
					ctrReg.Inc(1)

					if a, ok := p.Source.(source.Attribute); ok {
						return a.Sources, nil
					}

					return nil, nil
				},
			},
			"formats": &graphql.Field{
				Type:  graphql.NewNonNull( graphql.NewList( getAttributeFormatType() ) ),
				Description: "Scored Attributes Associated with the Country Model",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.format" )
					ctrReg.Inc(1)

					if a, ok := p.Source.(source.Attribute); ok {
						return a.Formats, nil
					}

					return nil, nil
				},
			},
		},
	})

	return attributeType
}

