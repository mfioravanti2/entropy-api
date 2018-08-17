package cli

import (
	"os"
	"strconv"
)

type Listener struct {
	Host string				`json:"host"`
	Port int				`json:"port"`
	UseTLS bool				`json:"use_TLS"`
	Encryption *Encryption	`json:"encryption,omitempty"`
}

func (l *Listener) EnvUpdate() {
	// Update the configuration options based on environment variable values
	if v := os.Getenv("ENTROPY_HOST"); v != "" {
		l.Host = v
	}
	if v := os.Getenv("ENTROPY_PORT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			l.Port = i
		}
	}
}

