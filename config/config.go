package config

import (
	"sync"
	"errors"
	"fmt"
)

const (
	MODE_SERVER string = "server"
	MODE_MIGRATE string = "migrate"
	MODE_EXAMPLE string = "save-config"
	MODE_UNKNOWN string = "unknown"

	DEFAULT_FILE string = "config.json"
	DEFAULT_HOST string = "127.0.0.1"
	DEFAULT_PORT int = 8080
)

type Config struct {
	modifyLock sync.RWMutex	`json:"-"`

	Config string			`json:"-"`

	// This should always be loaded from the command
	// line or from an environment variable
	// possible values: server, migrate, save
	Mode string				`json:"-"`

	Listener *Listener		`json:"listener,omitempty"`
	Security *Security		`json:"security,omitempty"`
	Logging *Logging		`json:"logging,omitempty"`
	Endpoints Endpoints		`json:"endpoints,omitempty"`
	Paths Paths				`json:"paths,omitempty"`
}

var config *Config

func init() {
	config, _ = DefaultConfig()
}

// Get the Global Configuration
func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	s := fmt.Sprintf("get config failed (no config)" )
	return nil, errors.New(s)
}

// Generate a Default Configuration
func DefaultConfig() (*Config, error) {
	c := &Config{}

	c.Listener = &Listener{}
	c.Listener.Host = DEFAULT_HOST
	c.Listener.Port = DEFAULT_PORT
	c.Listener.UseTLS = false

	c.Mode = MODE_UNKNOWN
	c.Config = DEFAULT_FILE

	c.Logging, _ = NewLogging()
	c.Endpoints, _ = NewEndpoints()
	c.Security, _ = NewSecurity()
	c.Paths, _ = NewPaths()

	return c, nil
}

// Set the Global Configuration
func SetConfig( c *Config ) error {
	if c != nil {
		config = c

		return nil
	}

	s := fmt.Sprintf("set config failed (no config)" )
	return errors.New(s)
}
