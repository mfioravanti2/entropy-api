package server

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc( func( w http.ResponseWriter, r *http.Request ) {
		rqId, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}

		httpContext := r.Context()
		reqCtx := logging.WithRqId(httpContext, rqId.String(), name, r.Method, r.RequestURI)

		logger := logging.Logger(reqCtx)
		logger.Info( "request received (logger)", )

		inner.ServeHTTP( w, r.WithContext( reqCtx ) )
	})
}
