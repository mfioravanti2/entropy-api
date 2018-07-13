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

func Run( c *cli.Config ) int {
	router := server.NewRouter()

	var connection string
	connection = fmt.Sprintf( "%s:%d", c.Host, c.Port )

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{ c.CorsOrigin })
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	corsRouter := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	var dataStore *scoringdb.DataStore
	var err error

	dataStore, err = scoringdb.GetDataStore( nil )
	if err == nil {
		defer dataStore.Close()

		switch c.Mode {
		case "server":
			if dataStore.Ready() {
				log.Fatal( http.ListenAndServe( connection, corsRouter ) )
			}
		case "migrate":
			dataStore.Migrate()
		default:
		}
	}

	return 0
}
