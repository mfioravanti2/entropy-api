package command

import (
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/handlers"

	"github.com/mfioravanti2/entropy-api/command/server"
	"github.com/mfioravanti2/entropy-api/cli"
)

func Run( c *cli.Config ) int {
	router := server.NewRouter()

	var connection string
	connection = fmt.Sprintf( "%s:%d", c.Host, c.Port )

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{ c.CorsOrigin })
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	corsRouter := handlers.CORS(originsOk, headersOk, methodsOk)(router)

	log.Fatal( http.ListenAndServe( connection, corsRouter ) )

	return 0
}
