package config

import (
	"context"
	"fmt"
	"errors"

	"github.com/google/uuid"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

const (
	ENDPOINT_REST    string = "entropy.api.policy.rest"
	ENDPOINT_GRAPHQL string = "entropy.api.policy.graphql"
	ENDPOINT_SCORING string = "entropy.api.policy.score"
	ENDPOINT_SYSTEM  string = "entropy.api.policy.system"
	ENDPOINT_METRICS string = "entropy.api.policy.metrics"
	ENDPOINT_HEALTH  string = "entropy.api.policy.health"
	ENDPOINT_OPENAPI string = "entropy.api.policy.openapi"
	ENDPOINT_DEFAULT string = "entropy.api.policy.default"
)

type Endpoint struct {
	Name string			`json:"name"`
	Enabled bool		`json:"enabled"`
	Restricted bool		`json:"restricted"`
	Entitlement string	`json:"entitlement"`
}

type Endpoints []Endpoint

func (e *Endpoint) DefaultConfig() {
	e.Name = ""
	e.Enabled = true
	e.Restricted = false
	e.Entitlement = ""

	entId, err := uuid.NewRandom()
	if err == nil {
		e.Entitlement = entId.String()
	}
}

// Create a default empty configuration
func NewEndpoint() (*Endpoint, error) {
	var e = &Endpoint{}
	e.DefaultConfig()

	ctx := logging.WithFuncId( context.Background(), "NewEndpoint", "config" )

	logger := logging.Logger( ctx )
	logger.Debug("generating default endpoint configuration",
	)

	return e, nil
}

func NewEndpoints() (Endpoints, error) {
	var eps = Endpoints{}
	var ep *Endpoint

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_REST
	eps = append( eps, *ep )

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_GRAPHQL
	eps = append( eps, *ep )

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_SCORING
	eps = append( eps, *ep )

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_METRICS
	eps = append( eps, *ep )

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_HEALTH
	eps = append( eps, *ep )

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_SYSTEM
	eps = append( eps, *ep )

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_OPENAPI
	eps = append( eps, *ep )

	ep, _ = NewEndpoint()
	ep.Name = ENDPOINT_DEFAULT
	ep.Enabled = false
	ep.Restricted = true
	eps = append( eps, *ep )

	return eps, nil
}

func (eps *Endpoints) GetEndpoint( name string ) (*Endpoint, error) {
	var ep Endpoint

	for _, ep = range *eps {
		if ep.Name == name {
			return &ep, nil
		}
	}

	s := fmt.Sprintf("find endpoint failed (endpoint name: %s)", name )
	return nil, errors.New(s)
}

