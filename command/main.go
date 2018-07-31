package command

import (
	"log"
	"net/http"
	"fmt"
	"github.com/gorilla/handlers"

	"github.com/mfioravanti2/entropy-api/command/server"
	"github.com/mfioravanti2/entropy-api/cli"
	"github.com/mfioravanti2/entropy-api/data/scoringdb"
)

// Exec the application based on the defined configuration
func Run( c *cli.Config ) int {
	router := server.NewRouter()

	var connection string
	connection = fmt.Sprintf( "%s:%d", c.Host, c.Port )

	// define the CORS headers, origins and methods
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{ c.CorsOrigin })
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// add the CORS handler
	corsRouter := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	var dataConfig *scoringdb.Config
	var dataStore *scoringdb.DataStore
	var err error

	// open the request/response scoring data store configuration
	// based on the application configuration
	dataConfig, err = scoringdb.OpenConfig( c.Files.DataStore )
	if err != nil {
		dataConfig, err = scoringdb.NewConfig()
	}

	// open the request/response scoring data store
	dataStore, err = scoringdb.GetDataStore( dataConfig )
	if err == nil {
		defer dataStore.Close()

		switch c.Mode {
		case "server":
			// start the service if running in server mode, but
			// only if the request/response data store is configured and available
			if dataStore.Ready( true ) {
				log.Fatal( http.ListenAndServe( connection, corsRouter ) )
			}
		case "migrate":
			// perform a data migration if the application is running in migrate mode
			dataStore.Migrate()
		default:
		}
	}

	return 0
}
