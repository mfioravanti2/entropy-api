package entropyql

import (
	"context"

	"go.uber.org/zap"
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func getAttributeFormatType() *graphql.Object {
	ctx := logging.WithFuncId( context.Background(), "getAttributeFormatType", "entropyql" )

	logger := logging.Logger( ctx )
	logger.Debug("building GraphQL schema",
		zap.String( "type", "attributeFormatType" ),
		)

	var formatType *graphql.Object

	formatType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Format",
		Description: "Calculation Format of a Scored Attribute",
		Fields: graphql.Fields{
			"format": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Method of Calculation for the Scored Value",
				Resolve: resolveFormatFormat,
			},
			"score": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Float),
				Description: "Scored Value of an Attribute",
				Resolve: resolveFormatScore,
			},
		},
	})

	return formatType
}

func getAttributeSourceType() *graphql.Object {
	ctx := logging.WithFuncId( context.Background(), "getAttributeSourceType", "entropyql" )

	logger := logging.Logger( ctx )
	logger.Debug("building GraphQL schema",
		zap.String( "type", "attributeSourceType" ),
		)

	var sourceType *graphql.Object

	sourceType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Source",
		Description: "Source of Frequency Statistics used to perform a Scoring Calculation",
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Document Title",
				Resolve: resolveSourceTitle,
			},
			"date": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
				Description: "Date of Publication or Retrieval",
				Resolve: resolveSourceTimestamp,
			},
			"organization": &graphql.Field{
				Type: graphql.String,
				Description: "Organization Responsible for the Source Document",
				Resolve: resolveSourcePublisher,
			},
			"url": &graphql.Field{
				Type: graphql.String,
				Description: "URL for retrieving the Source Document",
				Resolve: resolveSourceURL,
			},
		},
	})

	return sourceType
}

func getAttributeType() *graphql.Object {
	ctx := logging.WithFuncId( context.Background(), "getAttributeType", "entropyql" )

	logger := logging.Logger( ctx )
	logger.Debug("building GraphQL schema", zap.String( "type", "attributeType" ) )

	var attributeType *graphql.Object

	attributeType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Attribute",
		Description: "An attribute for a Country Model.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Globally unique identifier for the Attribute (UUID v4 format)",
				Resolve: resolveAttributeId,
			},
			"mnemonic": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Description: "Country Model unique identifier for the Attribute",
				Resolve: resolveAttributeMnemonic,
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Description: "Human friendly short description of the Attribute",
				Resolve: resolveAttributeName,
			},
			"notes": &graphql.Field{
				Type: graphql.String,
				Description: "Description of the Attribute",
				Resolve: resolveAttributeNotes,
			},
			"sources": &graphql.Field{
				Type: graphql.NewList( getAttributeSourceType() ),
				Description: "Source Documents for the Scoring Calculations",
				Resolve: resolveAttributeSources,
			},
			"formats": &graphql.Field{
				Type:  graphql.NewNonNull( graphql.NewList( getAttributeFormatType() ) ),
				Description: "Scored Attributes Associated with the Country Model",
				Args: graphql.FieldConfigArgument{
					"format": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: resolveAttributeFormats,
			},
		},
	})

	return attributeType
}

