package cli

import (
	"os"
	"strconv"
	"flag"
	"sync"
)

type Config struct {
	modifyLock sync.RWMutex

	Host string
	Port int

	CorsOrigin string

	Error error
}

func DefaultConfig() *Config {
	c := &Config{ Host: "localhost", Port: 8080, CorsOrigin: "*" }

	if err := c.ReadEnvironment(); err != nil {
		c.Error = err
		return c
	}

	return c
}

func (c *Config) ReadEnvironment() error {
	hostPtr := flag.String("host", "127.0.0.1", "Hostname")
	portPtr := flag.Int("port", 8080, "TCP port")
	corsPtr := flag.String("cors", "127.0.0.1", "CORS Origin Host")
	flag.Parse()

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

	c.modifyLock.Lock()
	defer c.modifyLock.Unlock()

	c.Host = *hostPtr
	c.Port = *portPtr
	c.CorsOrigin = *corsPtr

	return nil
}