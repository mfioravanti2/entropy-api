package command

import (
	"log"
	"net/http"
	"fmt"

	"github.com/mfioravanti2/entropy-api/command/server"
	"github.com/mfioravanti2/entropy-api/cli"
)

func Run( c *cli.Config ) int {
	router := server.NewRouter()

	var connection string
	connection = fmt.Sprintf( "%s:%d", c.Host, c.Port )

	log.Fatal(http.ListenAndServe(connection, router))

	return 0
}
