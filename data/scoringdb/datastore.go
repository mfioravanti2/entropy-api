package scoringdb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
	"github.com/mfioravanti2/entropy-api/data/scoringdb/none"
	"github.com/mfioravanti2/entropy-api/data/scoringdb/sqlite3"
	"github.com/mfioravanti2/entropy-api/data/scoringdb/mysql"
	"github.com/mfioravanti2/entropy-api/data/scoringdb/postgres"
)

type DataStore struct {
	Active	bool
	g		*gorm.DB
	c 		*Config
	LastUse time.Time
}

var dataStore *DataStore = nil

// Attempt to open a connection to the data store with the supplied configuration
func Open( c *Config ) ( *DataStore, error ) {
	// If the data store is already open, return the existing connection
	if dataStore != nil && dataStore.g != nil {
		return dataStore, nil
	}

	ctx := logging.WithFuncId( context.Background(), "Open", "scoringdb" )

	logger := logging.Logger( ctx )
	logger.Info("opening scoring data store object",
		zap.String( "engineId", c.Engine),
		zap.String( "connection", c.String() ),
	)

	var e error
	var g *gorm.DB
	var active bool

	// Open the connection based on the specified engine
	switch strings.ToLower( c.Engine ) {
	case none.ENGINE:
		g, active, e = none.Open( c.Connection )
	case sqlite3.ENGINE:
		g, active, e = sqlite3.Open( c.Connection )
	case mysql.ENGINE:
		g, active, e = mysql.Open( c.Connection )
	case postgres.ENGINE:
		g, active, e = postgres.Open( c.Connection )
	default:
		// unknown database engine was specified, specify that an error has occurred
		s := fmt.Sprintf("unknown database engine (%s)", c.Engine )
		e = errors.New( s )
	}

	if e != nil {
		// if any errors occurred, return the error an abort opening any connections
		logger.Error("error opening scoring data store object",
			zap.String( "engineId", c.Engine),
			zap.String( "connection", c.String() ),
			zap.String( "status", "error" ),
			zap.String("error ", e.Error() ),
		)

		return nil, e
	}

	logger.Info("assigning scoring data store configuration",
	)

	dataStore = &DataStore{ g: g, Active: active, c: c, LastUse: time.Now().UTC() }
	return dataStore, nil
}

// Perform a data store migration
func (ds *DataStore) Migrate() {
	// if the data store is not active, there is nothing to configure
	if !ds.Active {
		return
	}

	if ds != nil && ds.g != nil {
		var err error

		ctx := logging.WithFuncId( context.Background(), "Migrate", "scoringdb" )

		logger := logging.Logger( ctx )
		logger.Info("beginning scoring data store migration",
		)

		err = ds.g.AutoMigrate(&ReqRecord{}).Error
		if err != nil {
			logger.Error("scoring data store migration",
				zap.String( "recordType", "request"),
				zap.String( "tableId", "ReqRecord"),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}

		err = ds.g.AutoMigrate(&ReqAttribute{}).Error
		if err != nil {
			logger.Error("scoring data store migration",
				zap.String( "recordType", "request"),
				zap.String( "tableId", "ReqAttribute"),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}

		err = ds.g.AutoMigrate(&RespRecord{}).Error
		if err != nil {
			logger.Error("scoring data store migration",
				zap.String( "recordType", "response"),
				zap.String( "tableId", "RespRecord"),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}

		err = ds.g.AutoMigrate(&RespAttribute{}).Error
		if err != nil {
			logger.Error("scoring data store migration",
				zap.String( "recordType", "response"),
				zap.String( "tableId", "RespAttribute"),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}

		ds.LastUse = time.Now().UTC()

		logger.Info("completed scoring data store migration",
		)
	}
}

// Determine if the data store is ready to process requests
func (ds *DataStore) Ready( reqTables bool ) bool {
	if !ds.Active {
		return true
	}

	if err := ds.g.DB().Ping(); err == nil {
		ds.LastUse = time.Now().UTC()

		if reqTables {
			if ds != nil && ds.g != nil {
				return ds.readyRequest() && ds.readyResponse()
			}
		} else {
			return true
		}
	} else {
		ctx := logging.WithFuncId( context.Background(), "Ready", "scoringdb" )

		logger := logging.Logger( ctx )
		logger.Error("checking data store readiness",
			zap.String( "status", "error" ),
			zap.String("error ", err.Error() ),
		)
	}

	return false
}

// Close the connection to the data store
func (ds *DataStore) Close() {
	// if the data store is not active, there is no need to close it
	if !ds.Active {
		return
	}

	if ds != nil && ds.g != nil {
		ds.g.Close()

		ds.LastUse = time.Now().UTC()
		ds.g = nil

		ctx := logging.WithFuncId( context.Background(), "Close", "scoringdb" )

		logger := logging.Logger( ctx )
		logger.Info("closing scoring data store object",
		)
	}
}

// Return the configuration that was used to create the data store
func (ds *DataStore) Config() *Config  {
	return ds.c
}

// Get the active data store or generate a new data store based on
// the specified configuration
func GetDataStore( dbConfig *Config ) (*DataStore, error) {
	// if not data store is open, attempt to open one with the supplied configuration
	if dataStore == nil {
		var err error
		var config *Config

		// if no configuration was supplied generate the default configuration
		if dbConfig == nil {
			config, err = NewConfig()
			if err != nil {
				return nil, err
			}
		} else {
			config = dbConfig
		}

		// attempt to open the default configuration
		ds, err := Open( config )
		if err == nil {
			return ds, nil
		} else {
			return nil, err
		}
	}

	// return the existing data store
	return dataStore, nil
}
