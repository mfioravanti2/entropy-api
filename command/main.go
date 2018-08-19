package command

import (
	"fmt"
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/config/cli"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

// Exec the application based on the defined configuration
func Run( c *cli.CLI ) int {
	ctx := logging.WithFuncId( context.Background(), "Run", "command" )

	logger := logging.Logger( ctx )
	logger.Info("executing CLI command",
		zap.String( "mode", c.Mode ),
	)

	var err error
	var cfg *config.Config

	if c != nil {
		cfg, err = c.LoadConfig()
		if err != nil {
			panic( err )
		}

		config.SetConfig( cfg )
	}

	switch c.Mode {
	case config.MODE_SERVER:
		runServer( cfg )
	case config.MODE_MIGRATE:
		doMigrate( cfg )
	case config.MODE_EXAMPLE:
		saveExample( c )
	default:
		s := fmt.Sprintf("unknown execution mode (%s)", c.Mode )
		panic( errors.New(s) )
	}

	return 0
}

