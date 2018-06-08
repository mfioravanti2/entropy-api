package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"github.com/google/uuid"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

var httpContext = context.Background()

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc( func( w http.ResponseWriter, r *http.Request){
		rqId, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}

		rqCtx := logging.WithRqId(httpContext, rqId.String())
		logger := logging.Logger(rqCtx)

		logger.Info(name,
			zap.String("method", r.Method),
			zap.String("url", r.RequestURI),
		)

		inner.ServeHTTP(w,r)
	})
}
