package cli

import (
	"flag"
	"sync"
	"strings"
	"errors"
	"fmt"
)

const (
	MODE_SERVER string = "server"
	MODE_MIGRATE string = "migrate"
	MODE_EXAMPLE string = "save-config"
	MODE_UNKNOWN string = "unknown"
)

type Config struct {
	modifyLock sync.RWMutex	`json:"-"`

	Config string			`json:"-"`

	// This should always be loaded from the command line
	// possible values include: server, migrate, save
	Mode string				`json:"-"`

	Listener *Listener		`json:"listener,omitempty"`
	Security *Security		`json:"security,omitempty"`
	Logging *Logging		`json:"logging,omitempty"`
	Endpoints Endpoints		`json:"endpoints,omitempty"`
	Paths Paths				`json:"paths,omitempty"`
}

var config *Config

func init() {
	config, _ = DefaultConfig( true )
}

func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	s := fmt.Sprintf("get config failed (no config)" )
	return nil, errors.New(s)
}

// Generate an Environment Configuration
func DefaultConfig( useEnv bool ) (*Config, error) {
	c := &Config{}

	c.Listener = &Listener{}
	c.Listener.Host = "127.0.0.1"
	c.Listener.Port = 8080
	c.Listener.UseTLS = false

	c.Mode = MODE_UNKNOWN
	c.Config = "config.json"

	c.Logging, _ = NewLogging()
	c.Endpoints, _ = NewEndpoints()
	c.Security, _ = NewSecurity()
	c.Paths, _ = NewPaths()

	if useEnv {
		if err := c.ReadEnvironment(); err != nil {
			return c, err
		}
	}

	return c, nil
}

// Modify a Configuration by Over-riding Values with Command
// Line Parameters and Environment Variables
func (c *Config) ReadEnvironment() error {
	// Read configuration options from the command line parameters
	hostPtr := flag.String("host", "", "Hostname")
	portPtr := flag.Int("port", -1, "TCP port")
	cfgPtr := flag.String("config", "config.json", "Configuration File")
	modePtr := flag.String("mode", MODE_UNKNOWN, "Application Mode")
	flag.Parse()

	c.modifyLock.Lock()
	defer c.modifyLock.Unlock()

	c.Listener.EnvUpdate()

	c.Listener.Host = checkString( c.Listener.Host, *hostPtr, "" )
	c.Listener.Port = checkInt( c.Listener.Port, *portPtr, -1 )
	c.Mode =  checkString( c.Mode, strings.ToLower( *modePtr ), MODE_UNKNOWN )

	c.Config = *cfgPtr

	return nil
}

func checkString( currentValue string, flagValue string, defaultValue string ) string {
	if flagValue != defaultValue {
		return flagValue
	}

	return currentValue
}

func checkInt( currentValue, flagValue int, defaultValue int ) int {
	if flagValue != defaultValue {
		return flagValue
	}

	return currentValue
}