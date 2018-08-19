package command

import (
	"log"
	"fmt"
	"context"
	"net/http"

	"go.uber.org/zap"
	"github.com/gorilla/handlers"

	"github.com/mfioravanti2/entropy-api/data/scoringdb"
	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/command/server"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

func runServer( c *config.Config ) {
	ctx := logging.WithFuncId(context.Background(), "runServer", "command")

	logger := logging.Logger(ctx)
	logger.Info("running server",
		zap.String("host", c.Listener.Host),
		zap.Int("port", c.Listener.Port),
		zap.String("mode", c.Mode),
		zap.String("file", c.Config ),
	)

	router := server.NewRouter()

	var connection string
	connection = fmt.Sprintf( "%s:%d", c.Listener.Host, c.Listener.Port )

	// define the CORS headers, origins and methods
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// add the CORS handler
	corsRouter := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	var dataConfig *config.Backend
	var dataStore *scoringdb.DataStore
	var err error

	// open the request/response scoring data store configuration
	// based on the application configuration
	dataConfig = c.Logging.Backend

	// open the request/response scoring data store
	dataStore, err = scoringdb.GetDataStore(dataConfig)
	if err == nil {
		defer dataStore.Close()

		// start the service if running in server mode, but
		// only if the request/response data store is configured and available
		if dataStore.Ready(true) {
			logger.Debug("starting server",
				zap.String("host", c.Listener.Host),
				zap.Int("port", c.Listener.Port),
			)

			log.Fatal(http.ListenAndServe(connection, corsRouter))
		}
	}
}

