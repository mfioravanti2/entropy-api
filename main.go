package main

import (
	"github.com/mfioravanti2/entropy-api/command"
	"github.com/mfioravanti2/entropy-api/cli"
)

func main() {
	config := cli.DefaultConfig()
	command.Run( config )
}

