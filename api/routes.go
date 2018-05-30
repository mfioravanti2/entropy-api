package api

import (
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/api/attribute"
	"github.com/mfioravanti2/entropy-api/api/country"
	"github.com/mfioravanti2/entropy-api/api/sys"
	"github.com/mfioravanti2/entropy-api/api/scores"
)

func newRoutes() model.Routes {
	var routes = model.Routes{}

	routes = country.AddHandlers( routes )
	routes = attribute.AddHandlers( routes )
	routes = scores.AddHandlers( routes )
	routes = sys.AddHandlers( routes )

	return routes
}

func GetRoutes() model.Routes {
	return newRoutes()
}


