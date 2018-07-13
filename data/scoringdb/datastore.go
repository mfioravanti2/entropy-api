package scoringdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/data/scoringdb/sqlite3"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type DataStore struct {
	g *gorm.DB
}

var dataStore *DataStore = nil

func Open( c *Config ) ( *DataStore, error ) {
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

	switch c.Engine {
	case sqlite3.ENGINE:
		g, e = sqlite3.Open( c.Connection )
		if e != nil {
			logger.Error("error opening scoring data store object",
				zap.String( "engineId", c.Engine),
				zap.String( "connection", c.String() ),
				zap.String( "status", "error" ),
				zap.String("error ", e.Error() ),
			)

			return nil, e
		}
	default:
		s := fmt.Sprintf("unknown database engine (%s)", c.Engine )

		logger.Error("error opening scoring data store object",
			zap.String( "engineId", c.Engine),
			zap.String( "connection", c.String() ),
			zap.String( "status", "error" ),
			zap.String("error ", s ),
		)

		return nil, errors.New(s)
	}

	logger.Info("assiging scoring data store configuration",
	)

	dataStore = &DataStore{ g: g }
	return dataStore, nil
}

func (ds *DataStore) Migrate() {
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

		logger.Info("completed scoring data store migration",
		)
	}
}

func (ds *DataStore) Ready() bool {
	if ds != nil && ds.g != nil {
		return ds.readyRequest() && ds.readyResponse()
	}

	return false
}

func (ds *DataStore) Close() {
	if ds != nil && ds.g != nil {
		ds.g.Close()

		ds.g = nil

		ctx := logging.WithFuncId( context.Background(), "Close", "scoringdb" )

		logger := logging.Logger( ctx )
		logger.Info("closing scoring data store object",
		)
	}
}

func GetDataStore( dbConfig *Config ) (*DataStore, error) {
	if dataStore == nil {
		var err error
		var config *Config

		if dbConfig == nil {
			config, err = NewConfig()
			if err != nil {
				return nil, err
			}
		} else {
			config = dbConfig
		}

		ds, err := Open( config )
		if err == nil {
			return ds, nil
		} else {
			return nil, err
		}
	}

	return dataStore, nil
}
