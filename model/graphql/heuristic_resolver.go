package entropyql

import (
	"github.com/graphql-go/graphql"

	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/model/source"
)

func resolveHeuristicId(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.heuristics.id" )
	ctrReg.Inc(1)

	if model, ok := p.Source.(source.Heuristic); ok {
		return model.Id, nil
	}

	return nil, nil
}

func resolveHeuristicNotes(p graphql.ResolveParams) (interface{}, error) {
	if model, ok := p.Source.(source.Heuristic); ok {
		return model.Notes, nil
	}

	return nil, nil
}

func resolveHeuristicMatch(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.heuristics.match" )
	ctrReg.Inc(1)

	if model, ok := p.Source.(source.Heuristic); ok {
		return model.Match, nil
	}

	return nil, nil
}

func resolveHeuristicInsert(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.heuristics.insert" )
	ctrReg.Inc(1)

	if model, ok := p.Source.(source.Heuristic); ok {
		return model.Insert, nil
	}

	return nil, nil
}

func resolveHeuristicRemove(p graphql.ResolveParams) (interface{}, error) {
	ctrReg, _ := metrix.GetCounter( "entropy.graphql.query.heuristics.remove" )
	ctrReg.Inc(1)

	if model, ok := p.Source.(source.Heuristic); ok {
		return model.Remove, nil
	}

	return nil, nil
}
