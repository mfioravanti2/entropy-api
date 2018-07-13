package scoringdb

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/mfioravanti2/entropy-api/data/scoringdb/sqlite3"
)

type DataStore struct {
	g *gorm.DB
}

var dataStore *DataStore = nil

func Open( c *Config ) ( *DataStore, error ) {
	if dataStore != nil {
		return dataStore, nil
	}

	var e error
	var g *gorm.DB

	switch c.Engine {
	case sqlite3.ENGINE:
		g, e = sqlite3.Open( c.Connection )
		if e != nil {
			return nil, e
		}
	default:
		s := fmt.Sprintf("unknown database engine (%s)", c.Engine )
		return nil, errors.New(s)
	}

	dataStore = &DataStore{ g: g }
	return dataStore, nil
}

func (ds *DataStore) Migrate() {
	if ds != nil && ds.g != nil {
		ds.g.AutoMigrate(&ReqRecord{})
		ds.g.AutoMigrate(&ReqAttribute{})
	}
}

func (ds *DataStore) Ready() bool {
	if ds != nil && ds.g != nil {
		if ds.g.HasTable(&ReqRecord{}) {
			if ds.g.HasTable(&ReqAttribute{}) {
				return true
			}
		}
	}

	return false
}

func (ds *DataStore) Close() {
	if ds != nil && ds.g != nil {
		ds.g.Close()
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
