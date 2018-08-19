package cli

import (
	"fmt"
	"flag"
	"errors"
	"sync"
	"strings"
	"os"
	"strconv"

	"github.com/mfioravanti2/entropy-api/config"
	"io/ioutil"
	"encoding/json"
)

const (
	ENV_VAR_FILE string = "ENTROPY_CONFIG"
	ENV_VAR_HOST string = "ENTROPY_HOST"
	ENV_VAR_PORT string = "ENTROPY_PORT"
	ENV_VAR_MODE string = "ENTROPY_MODE"
)

type CLI struct {
	modifyLock sync.RWMutex	`json:"-"`

	File string				`json:"-"`
	Mode string				`json:"-"`

	Host string				`json:"-"`
	Port int				`json:"-"`
}

var cli *CLI

func init() {
	cli = &CLI{ File: config.DEFAULT_FILE,
				Mode: config.MODE_UNKNOWN,
				Host: config.DEFAULT_HOST,
				Port: config.DEFAULT_PORT }
}

func GetCLI() (*CLI, error) {
	if cli != nil {
		return cli, nil

	}

	s := fmt.Sprintf("get command line arguments failed (no object)" )
	return nil, errors.New(s)
}

func (c *CLI) ReadArgs() {
	hostPtr := flag.String("host", config.DEFAULT_HOST, "Host IP Address")
	portPtr := flag.Int("port", config.DEFAULT_PORT, "TCP port")
	filePtr := flag.String("config", config.DEFAULT_FILE, "Configuration File")
	modePtr := flag.String("mode", config.MODE_UNKNOWN, "Application Mode")
	flag.Parse()

	c.modifyLock.Lock()
	defer c.modifyLock.Unlock()

	c.Host = checkString( c.Host, *hostPtr, config.DEFAULT_HOST )
	c.Port = checkInt( c.Port, *portPtr, config.DEFAULT_PORT )
	c.Mode =  checkString( c.Mode, strings.ToLower( *modePtr ), config.MODE_UNKNOWN )
	c.File =  checkString( c.Mode, *filePtr, config.MODE_UNKNOWN )
}

// Update the configuration options based on environment variable values
func (c *CLI) EnvUpdate() {
	c.modifyLock.Lock()
	defer c.modifyLock.Unlock()

	// Check to see if a host ip has been specified as an environment variable
	if v := os.Getenv(ENV_VAR_HOST); v != "" {
		// If the host is not the default and the CLI host is
		// the default, use the environment's setting
		if v != config.DEFAULT_HOST && c.Host == config.DEFAULT_HOST {
			c.Host = v
		}
	}

	// Check to see if a port number has been specified as an environment variable
	if v := os.Getenv(ENV_VAR_PORT); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			// If the port is not the default and the CLI port is
			// the default, use the environment's setting
			if i != config.DEFAULT_PORT && c.Port == config.DEFAULT_PORT {
				c.Port = i
			}
		}
	}

	// Check to see if a mode has been specified as an environment variable
	if v := os.Getenv(ENV_VAR_MODE); v != "" {
		// If the mode is not the default and the CLI mode is
		// the default, use the environment's setting
		if v != config.MODE_UNKNOWN && c.Mode == config.MODE_UNKNOWN {
			c.Mode = v
		}
	}

	// Check to see if a configuration file has been specified as an environment variable
	if v := os.Getenv(ENV_VAR_FILE); v != "" {
		// If the configuration file name is not the default and the CLI configuration
		// file name is the default, use the environment's setting
		if v != config.DEFAULT_FILE && c.File == config.DEFAULT_FILE {
			c.File = v
		}
	}
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

func ( c *CLI ) LoadConfig() ( *config.Config, error ) {
	var cfg config.Config

	if _, err := os.Stat( c.File ); os.IsNotExist(err) {
		s := fmt.Sprintf("configuration file not found (%s)", c.File )
		return nil, errors.New(s)
	}

	jsonData, err := ioutil.ReadFile( c.File )
	if err != nil {
		s := fmt.Sprintf("unable to load configuration file (%s)", c.File )
		return nil, errors.New(s)
	}

	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		s := fmt.Sprintf("unable to parse source file (%s), expected json format", c.File )
		return nil, errors.New(s)
	}

	cfg.Mode = c.Mode
	cfg.Config = c.File

	cfg.Listener.Host = c.Host
	cfg.Listener.Port = c.Port

	return &cfg, nil
}
