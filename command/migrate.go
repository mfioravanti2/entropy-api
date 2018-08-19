package command

import (
	"context"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/data/scoringdb"
)

func doMigrate( c *config.Config ) {
	ctx := logging.WithFuncId( context.Background(), "doMigrate", "command" )

	logger := logging.Logger( ctx )
	logger.Info("beginning database migration",
		zap.String( "mode", c.Mode ),
		zap.String( "file", c.Config ),
	)

	var dataConfig *config.Backend
	var dataStore *scoringdb.DataStore
	var err error

	// open the request/response scoring data store
	dataStore, err = scoringdb.GetDataStore( dataConfig )
	if err == nil {
		defer dataStore.Close()

		dataStore.Migrate()
	} else {
		logger.Error("unable to access database configuration",
			zap.String( "mode", c.Mode ),
			zap.String( "file", c.Config ),
			zap.String( "status", "error"),
			zap.String( "error", err.Error() ),
		)
	}
}
