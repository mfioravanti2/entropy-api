package api

import (
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/api/attribute"
	"github.com/mfioravanti2/entropy-api/api/country"
)

func newRoutes() model.Routes {
	var routes = model.Routes{}

	routes = country.AddHandlers( routes )
	routes = attribute.AddHandlers( routes )

	return routes
}

func GetRoutes() model.Routes {
	return newRoutes()
}


