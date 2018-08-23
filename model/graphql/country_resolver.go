package entropyql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model/source"
	"github.com/mfioravanti2/entropy-api/model/metrics"
)

func resolveCountryLocale(p graphql.ResolveParams) (interface{}, error) {
	if model, ok := p.Source.(source.Model); ok {
		return model.Locale, nil
	}

	return nil, nil
}

func resolveCountryThreshold(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.country.threshold" )
	ctrReg.Inc(1)

	if model, ok := p.Source.(source.Model); ok {
		return model.Threshold, nil
	}

	return nil, nil
}

func resolveCountryK(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.country.k" )
	ctrReg.Inc(1)

	if model, ok := p.Source.(source.Model); ok {
		return model.K, nil
	}

	return nil, nil
}

func resolveCountryVersion(p graphql.ResolveParams) (interface{}, error) {
	if model, ok := p.Source.(source.Model); ok {
		return model.ModelVersion, nil
	}

	return nil, nil
}

func resolveCountryTimestamp(p graphql.ResolveParams) (interface{}, error) {
	if model, ok := p.Source.(source.Model); ok {
		return model.ModelDate, nil
	}

	return nil, nil
}

func resolveCountryHeuristics(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.heuristics" )
	ctrReg.Inc(1)

	heuristics := make([]source.Heuristic, 0)

	var ok bool
	var model source.Model

	if model, ok = p.Source.(source.Model); !ok {
		s := fmt.Sprintf( "unable to convert parameter source, expected source.Model" )
		return heuristics, errors.New(s)
	}

	if len(p.Args) == 0 {
		return model.Heuristics, nil
	} else {
		ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.heuristics.args" )
		ctrReg.Inc(1)

		var id string

		if id, ok = p.Args["id"].(string); !ok {
			s := fmt.Sprintf( "unable to parse parameter argument (id)" )
			return heuristics, errors.New(s)
		}

		for _, h := range model.Heuristics {
			if h.Id == id {
				heuristics = append( heuristics, h )
			}
		}

		return heuristics, nil
	}
}

func resolveCountryAttributes(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attributes" )
	ctrReg.Inc(1)

	attributes := make([]source.Attribute, 0)

	var ok bool
	var model source.Model

	if model, ok = p.Source.(source.Model); !ok {
		s := fmt.Sprintf( "unable to convert parameter source, expected source.Model" )
		return attributes, errors.New(s)
	}

	if len(p.Args) == 0 {
		return model.Attributes, nil
	} else {
		ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.attributes.args" )
		ctrReg.Inc(1)

		var id string

		if id, ok = p.Args["mnemonic"].(string); !ok {
			s := fmt.Sprintf( "unable to parse parameter argument (mnemonic)" )
			return attributes, errors.New(s)
		}

		for _, a := range model.Attributes {
			if a.Mnemonic == id {
				attributes = append( attributes, a )
			}
		}

		return attributes, nil
	}
}
