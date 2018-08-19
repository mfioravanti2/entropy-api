package model

import (
	"net/http"

	"github.com/mfioravanti2/entropy-api/config"
)

const (
	// Allow Content-Type enforcement on the request type
	ENFORCE_CONTENT_JSON	string = "json"
	ENFORCE_CONTENT_GRAPHQL	string = "graphql"
	ENFORCE_CONTENT_NONE	string = "none"

	// Specify the type of authN required for the route
	AUTH_METHOD_JWT			string = "jwt"
	AUTH_METHOD_NONE		string = "none"
)

// Endpoint Routing Definition
type Route struct {
	Name		string
	Method		string
	Pattern		string
	HandlerFunc	http.HandlerFunc
	Params		[]string

	// Content-type enforcement on the request
	Enforce		string
	// Endpoint Policy configuration
	Policy		*config.Endpoint
	// Authentication type enforcement (if enabled in the configuration)
	AuthN		string
}

type Routes []Route

