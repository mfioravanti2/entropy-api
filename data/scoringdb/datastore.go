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

	var g *gorm.DB
	var e error

	switch c.Engine {
	case sqlite3.ENGINE:
		g, e = sqlite3.Open()
	default:
		s := fmt.Sprintf("unknown database engine (%s)", c.Engine )
		e = errors.New(s)
		return nil, e
	}

	if !g.HasTable(&ReqRecord{}) {
		g.CreateTable(&ReqRecord{})
		g.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&ReqRecord{})

		if !g.HasTable(&ReqAttribute{}) {
			g.CreateTable(&ReqAttribute{})
			g.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&ReqAttribute{})
		}

		g.Model(&ReqAttribute{}).AddForeignKey("req_record_id", "records(req_record_id)", "CASCADE", "CASCADE")
	}

	dataStore = &DataStore{ g: g}

	return dataStore, e
}

func (ds *DataStore) Close() {
	ds.g.Close()
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
		if err != err {
			return ds, nil
		} else {
			return nil, err
		}
	}

	return dataStore, nil
}
