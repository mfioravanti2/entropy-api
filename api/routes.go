/*
	Generate an array which contains all of the REST API endpoint connections
*/
package api

import (
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/config"

	"github.com/mfioravanti2/entropy-api/api/attribute"
	"github.com/mfioravanti2/entropy-api/api/country"
	"github.com/mfioravanti2/entropy-api/api/sys"
	"github.com/mfioravanti2/entropy-api/api/scores"
	"github.com/mfioravanti2/entropy-api/api/heuristic"
	"github.com/mfioravanti2/entropy-api/api/openapi-spec"
	"github.com/mfioravanti2/entropy-api/api/metrics"
	"github.com/mfioravanti2/entropy-api/api/graphql"
)

//	Generate a complete list of available routes
func newRoutes() model.Routes {
	var routes = model.Routes{}
	var c *config.Config

	c, err := config.GetConfig()
	if err == nil {
		routes = country.AddHandlers( routes, &c.Endpoints )
		routes = attribute.AddHandlers( routes, &c.Endpoints )
		routes = heuristic.AddHandlers( routes, &c.Endpoints )
		routes = scores.AddHandlers( routes, &c.Endpoints )
		routes = sys.AddHandlers( routes, &c.Endpoints )
		routes = openapi_spec.AddHandlers( routes, &c.Endpoints )
		routes = metrics.AddHandlers( routes, &c.Endpoints )
		routes = graphql.AddHandlers( routes, &c.Endpoints )
	}

	return routes
}

//	Get a complete list of available routes
func GetRoutes() model.Routes {
	return newRoutes()
}


