package cli

import (
	"os"
	"strconv"
	"flag"
	"sync"
	"strings"
)

type Locations struct {
	DataStore	string
	Models		string
}

type Config struct {
	modifyLock sync.RWMutex

	Host string
	Port int

	CorsOrigin string

	Mode string

	Files Locations

	Error error
}

// Generate an Environment Configuration
func DefaultConfig() *Config {
	c := &Config{ Host: "localhost", Port: 8080, CorsOrigin: "*" }

	if err := c.ReadEnvironment(); err != nil {
		c.Error = err
		return c
	}

	return c
}

// Modify a Configuration by Over-riding Values with Command
// Line Parameters and Environment Variables
func (c *Config) ReadEnvironment() error {
	// Read configuration options from the command line parameters
	hostPtr := flag.String("host", "127.0.0.1", "Hostname")
	portPtr := flag.Int("port", 8080, "TCP port")
	corsPtr := flag.String("cors", "127.0.0.1", "CORS Origin Host")
	modePtr := flag.String("mode", "server", "Application Mode")
	dbCfPtr := flag.String("db_config", "test/sqlite_config.json", "DataStore Config")
	modlPtr := flag.String("models", "data/sources/", "Country Model Directory")
	flag.Parse()

	// Update the configuration options based on environment variable values
	if v := os.Getenv("ENTROPY_HOST"); v != "" {
		hostPtr = &v
	}
	if v := os.Getenv("ENTROPY_PORT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			portPtr = &i
		} else {
			panic(err)
		}
	}
	if v := os.Getenv("ORIGIN_ALLOWED"); v != "" {
		corsPtr = &v
	}
	if v := os.Getenv("MODEL_PATH"); v != "" {
		modlPtr = &v
	}

	c.modifyLock.Lock()
	defer c.modifyLock.Unlock()

	c.Host = *hostPtr
	c.Port = *portPtr
	c.CorsOrigin = *corsPtr
	c.Mode =  strings.ToLower( *modePtr )
	c.Files.DataStore = *dbCfPtr

	if strings.HasSuffix( *modlPtr, "/" ) {
		c.Files.Models = *modlPtr
	} else {
		c.Files.Models = ( *modlPtr ) + "/"
	}

	return nil
}