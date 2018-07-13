package scoringdb

import (
	"context"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"errors"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type Config struct {
	Engine		string	`json:"engine"`
	Connection	string	`json:"connection"`
	Hide		bool	`json:"hide"`
}

func (c *Config) String() string {
	if c.Hide {
		return "CONNECTION STRING HIDDEN"
	}

	return c.Connection
}

func NewConfig() (*Config, error) {
	c := Config{ Engine: "sqlite3", Connection: "./scores.db", Hide: false }

	ctx := logging.WithFuncId( context.Background(), "NewConfig", "scoringdb" )

	logger := logging.Logger( ctx )
	logger.Info("generating default scoring data store configuration",
		zap.String("engineId", c.Engine ),
		zap.String("connection", c.String() ),
	)

	return &c, nil
}

func OpenConfig( configFile string ) ( *Config, error ) {
	ctx := logging.WithFuncId( context.Background(), "OpenConfig", "scoringdb" )

	logger := logging.Logger( ctx )
	logger.Info("preparing to load scoring data store configuration",
		zap.String("config_file", configFile ),
	)

	jsonData, err := ioutil.ReadFile( configFile )
	if err != nil {
		s := fmt.Sprintf("unable to load configuration file")
		logger.Error( "loading configuration file",
			zap.String("file", configFile ),
			zap.String("error", s ),
		)

		return nil, errors.New(s)
	}

	logger.Debug("attempting to unmarshal data score configuration",
		zap.String("config_file", configFile ),
	)

	var c Config
	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		s := fmt.Sprintf("unable to parse configuration file, expected json format")
		logger.Error( "unable to parse configuration file",
			zap.String("file", configFile ),
			zap.String("error", s ),
		)

		return nil, errors.New(s)
	}

	logger.Info("loaded scoring data store configuration",
		zap.String("engineId", c.Engine ),
		zap.String("connection", c.String() ),
	)

	return &c, nil
}
