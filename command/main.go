package command

import (
	"log"
	"net/http"
	"fmt"

	"github.com/mfioravanti2/entropy-api/command/server"
)

func Run(host string, port int) int {
	router := server.NewRouter()

	var connection string
	connection = fmt.Sprintf( "%s:%d", host, port )

	log.Fatal(http.ListenAndServe(connection, router))

	return 0
}
