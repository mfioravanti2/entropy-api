package config

import (
	"context"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type Backend struct {
	Engine		string	`json:"engine"`
	Connection	string	`json:"connection"`
	Hide		bool	`json:"hide"`
	Redacted	string	`json:"redacted"`
}

func (b *Backend) DefaultConfig() {
	b.Engine = "none"
	b.Connection = ""
	b.Hide = false
	b.Redacted = ""
}

// Return the Connection string. If the Hide flag is set,
// the Redacted string will be logged
func (b *Backend) String() string {
	if b.Hide {
		return b.Redacted
	}

	return b.Connection
}

// Create a default empty configuration
func NewBackend() (*Backend, error) {
	var b = &Backend{}
	b.DefaultConfig()

	ctx := logging.WithFuncId( context.Background(), "NewBackend", "config" )

	logger := logging.Logger( ctx )
	logger.Debug("generating default backend configuration",
	)

	return b, nil
}

