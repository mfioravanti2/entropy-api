package entropyql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/model/source"
	"github.com/mfioravanti2/entropy-api/data"
)

func resolveCountries(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.countries" )
	ctrReg.Inc(1)

	countries := make([]source.Model, 0)

	if len(p.Args) == 0 {
		return data.GetAllCountries(), nil
	} else {
		ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.countries.args" )
		ctrReg.Inc(1)

		var ok bool
		var locale string

		if locale, ok = p.Args["locale"].(string); !ok {
			s := fmt.Sprintf( "unable to parse parameter argument (locale)" )
			return countries, errors.New(s)
		}

		for _, c := range data.GetAllCountries() {
			if c.Locale == locale {
				countries = append( countries, c )
			}
		}

		return countries, nil
	}
}
