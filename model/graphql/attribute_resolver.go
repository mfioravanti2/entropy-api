package entropyql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/model/source"
)

func resolveFormatFormat(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.format.method" )
	ctrReg.Inc(1)

	if f, ok := p.Source.(source.Format); ok {
		return f.Format, nil
	}

	return nil, nil
}

func resolveFormatScore(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.format.score" )
	ctrReg.Inc(1)

	if f, ok := p.Source.(source.Format); ok {
		return f.Score, nil
	}

	return nil, nil
}

func resolveSourceTitle(p graphql.ResolveParams) (interface{}, error) {
	if s, ok := p.Source.(source.Source); ok {
		return s.Title, nil
	}

	return nil, nil
}

func resolveSourceTimestamp(p graphql.ResolveParams) (interface{}, error) {
	if s, ok := p.Source.(source.Source); ok {
		return s.Date, nil
	}

	return nil, nil
}

func resolveSourcePublisher(p graphql.ResolveParams) (interface{}, error) {
	if s, ok := p.Source.(source.Source); ok {
		return s.Org, nil
	}

	return nil, nil
}

func resolveSourceURL(p graphql.ResolveParams) (interface{}, error) {
	if s, ok := p.Source.(source.Source); ok {
		return s.URI, nil
	}

	return nil, nil
}

func resolveAttributeId(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.id" )
	ctrReg.Inc(1)

	if a, ok := p.Source.(source.Attribute); ok {
		return a.Id, nil
	}

	return nil, nil
}

func resolveAttributeMnemonic(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.mnemonic" )
	ctrReg.Inc(1)

	if a, ok := p.Source.(source.Attribute); ok {
		return a.Mnemonic, nil
	}

	return nil, nil
}

func resolveAttributeName(p graphql.ResolveParams) (interface{}, error) {
	if a, ok := p.Source.(source.Attribute); ok {
		return a.Name, nil
	}

	return nil, nil
}

func resolveAttributeNotes(p graphql.ResolveParams) (interface{}, error) {
	if a, ok := p.Source.(source.Attribute); ok {
		return a.Notes, nil
	}

	return nil, nil
}

func resolveAttributeSources(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.sources" )
	ctrReg.Inc(1)

	if a, ok := p.Source.(source.Attribute); ok {
		return a.Sources, nil
	}

	return nil, nil
}

func resolveAttributeFormats(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.format" )
	ctrReg.Inc(1)

	formats := make([]source.Format, 0)

	var ok bool
	var attribute source.Attribute

	if attribute, ok = p.Source.(source.Attribute); !ok {
		s := fmt.Sprintf( "unable to convert parameter source, expected source.Attribute" )
		return formats, errors.New(s)
	}

	if len(p.Args) == 0 {
		return attribute.Formats, nil
	} else {
		ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attribute.format.args" )
		ctrReg.Inc(1)

		var format string

		if format, ok = p.Args["format"].(string); !ok {
			s := fmt.Sprintf( "unable to parse parameter argument (format)" )
			return formats, errors.New(s)
		}

		for _, f := range attribute.Formats {
			if f.Format == format {
				formats = append( formats, f )
			}
		}

		return formats, nil
	}
}