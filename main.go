/*

entropy-api is a command line to which provides a RESTful API which attempts
to assist in determining if a set of attributes contains enough information
that is should be treated as Personally Identifiable Information (or Personal
Data).

Usage:
	entropy-api [host="<Listening IP address>"]
				[-port="<TCP Listening port>"]
				[-mode='mode']
				[-config=/path/to/config]

	The commands are
	host	 IP address to listen for connections (only in server mode)
			 default: 127.0.0.1
	port	 TCP port to listen for connections (only in server mode)
			 default: 8080
	mode	 Application Mode (i.e. server, migrate, save-config)
			 server: Run the application in a model that listens on the HOST:PORT for connections
			 migrate: Connect to the Request/Response Data Store and setup the tables
	config	 Configuration file
			 if the mode is 'server' or 'migrate' the configuration will be used, otherwise if
			 the mode is 'save-config' the file will be overwritten with an example configuration

*/
package main

import (
	"github.com/mfioravanti2/entropy-api/command"
	"github.com/mfioravanti2/entropy-api/config/cli"
)

func main() {
	// Parse configuration options from command line args
	// and environment variables. Command line args take
	// precedence over environment variables
	c, err := cli.GetCLI()
	if err != nil {
		panic(err)
	}
	c.ReadArgs()
	c.EnvUpdate()

	command.Run( c )
}

