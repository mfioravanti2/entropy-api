package command

import (
	"context"
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/handlers"
	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/command/server"
	"github.com/mfioravanti2/entropy-api/cli"
	"github.com/mfioravanti2/entropy-api/data/scoringdb"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

// Exec the application based on the defined configuration
func Run( c *cli.Config ) int {
	router := server.NewRouter()

	ctx := logging.WithFuncId( context.Background(), "Run", "command" )

	logger := logging.Logger( ctx )
	logger.Info("executing cli command",
		zap.String( "host", c.Listener.Host ),
		zap.Int( "port", c.Listener.Port ),
		zap.String( "mode", c.Mode ),
	)

	var connection string
	connection = fmt.Sprintf( "%s:%d", c.Listener.Host, c.Listener.Port )

	// define the CORS headers, origins and methods
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// add the CORS handler
	corsRouter := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	var dataConfig *cli.Backend
	var dataStore *scoringdb.DataStore
	var err error

	// open the request/response scoring data store configuration
	// based on the application configuration
	dataConfig = c.Logging.Backend

	// open the request/response scoring data store
	dataStore, err = scoringdb.GetDataStore( dataConfig )
	if err == nil {
		defer dataStore.Close()

		switch c.Mode {
		case cli.MODE_SERVER:
			logger.Info("preparing to start server",
			)

			// start the service if running in server mode, but
			// only if the request/response data store is configured and available
			if dataStore.Ready( true ) {
				logger.Debug("starting server",
					zap.String( "host", c.Listener.Host ),
					zap.Int( "port", c.Listener.Port ),
				)

				log.Fatal( http.ListenAndServe( connection, corsRouter ) )
			}
		case cli.MODE_MIGRATE:
			// perform a data migration if the application is running in migrate mode
			logger.Info("initiating database migration",
			)

			dataStore.Migrate()
		case cli.MODE_EXAMPLE:
			logger.Info("starting saving sample configuration",
				zap.String( "file", c.Config ),
			)

			var cfg *cli.Config
			var err error

			cfg, err = cli.DefaultConfig( false )
			cfgJson, err := json.Marshal(cfg)
			err = ioutil.WriteFile( c.Config, cfgJson, 0644)
			if err != nil {
				panic(err)
			}
		default:
		}
	}

	return 0
}
