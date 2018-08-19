package command

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/config/cli"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

// Save an example configuration file, the example will be
// written to the file specified in CLI.File
func saveExample( c *cli.CLI ) {
	ctx := logging.WithFuncId( context.Background(), "saveExample", "command" )

	logger := logging.Logger( ctx )
	logger.Info("writing example configuration file",
		zap.String( "mode", c.Mode ),
		zap.String( "file", c.File ),
	)

	var cfg *config.Config
	var err error

	cfg, err = config.DefaultConfig()
	cfgJson, err := json.Marshal(cfg)
	err = ioutil.WriteFile( c.File, cfgJson, 0644)
	if err != nil {
		panic(err)
	}
}