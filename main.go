package main

import (
	"flag"

	"github.com/mfioravanti2/entropy-api/command"
)

func main() {
	hostPtr := flag.String("host", "127.0.0.1", "Hostname")
	portPtr := flag.Int("port", 8080, "TCP port")
	flag.Parse()

	command.Run( *hostPtr, *portPtr )
}

