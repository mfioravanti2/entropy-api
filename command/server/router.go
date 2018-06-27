package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mfioravanti2/entropy-api/api"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash( true )
	for _, route := range api.GetRoutes() {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		if route.Params == nil {
			router.
				Methods( route.Method ).
				Path( route.Pattern ).
				Name( route.Name ).
				Handler( handler )
		} else {
			router.
				Methods( route.Method ).
				Path( route.Pattern ).
				Queries( route.Params... ).
				Name( route.Name ).
				Handler( handler )
		}
	}

	return router
}

