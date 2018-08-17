package cli

import (
	"context"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type Logging struct {
	Mode string			`json:"mode"`
	Requests bool		`json:"requests"`
	Responses bool		`json:"responses"`
	Backend *Backend	`json:"backend"`
}

func (l *Logging) DefaultConfig() {
	l.Mode = "production"
	l.Requests = true
	l.Responses = true

	var b *Backend
	var err error

	b, err = NewBackend()
	if err == nil {
		l.Backend = b
	}
}

// Create a default empty configuration
func NewLogging() (*Logging, error) {
	var l = &Logging{}
	l.DefaultConfig()

	ctx := logging.WithFuncId( context.Background(), "NewLogging", "cli" )

	logger := logging.Logger( ctx )
	logger.Debug("generating default logging configuration",
	)

	return l, nil
}
