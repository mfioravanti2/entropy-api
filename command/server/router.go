package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/api"
	"github.com/mfioravanti2/entropy-api/command/server/headers"
	"github.com/mfioravanti2/entropy-api/command/server/enforce"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash( true )
	for _, route := range api.GetRoutes() {
		var handler http.Handler

		handler = route.HandlerFunc

		switch route.Enforce {
		case model.ENFORCE_CONTENT_JSON:
			handler = enforce.EnforceJSONHandler( handler )
		case model.ENFORCE_CONTENT_GRAPHQL:
		case model.ENFORCE_CONTENT_NONE:
		default:
		}

		switch route.AuthN {
		case model.AUTH_METHOD_JWT:
		case model.AUTH_METHOD_NONE:
		default:
		}

		handler = Logger( headers.SecurityHeadersHandler( handler ), route.Name )

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

