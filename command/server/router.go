package server

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/mfioravanti2/entropy-api/api"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash( true )
	for _, route := range api.GetRoutes() {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

