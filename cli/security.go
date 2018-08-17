package cli

import (
	"context"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type Encryption struct {
	Certificate string	`json:"cert"`
	PrivateKey string	`json:"key"`
}

type Security struct {
	Headers Headers			`json:"headers,omitempty"`
	Encryption *Encryption	`json:"encryption,omitempty"`
}

func (s *Security) DefaultConfig() {
	s.Encryption, _ = NewEncryption()
	s.Headers, _ = NewHeaders()
}

// Create a default empty configuration
func NewSecurity() (*Security, error) {
	var s = &Security{}
	s.DefaultConfig()

	ctx := logging.WithFuncId( context.Background(), "NewSecurity", "cli" )

	logger := logging.Logger( ctx )
	logger.Debug("generating default logging configuration",
	)

	return s, nil
}

func (e *Encryption) DefaultConfig() {
	e.Certificate = "server.pem"
	e.PrivateKey = "server-key.pem"
}

// Create a default empty configuration
func NewEncryption() (*Encryption, error) {
	var e = &Encryption{}
	e.DefaultConfig()

	ctx := logging.WithFuncId( context.Background(), "NewEncryption", "cli" )

	logger := logging.Logger( ctx )
	logger.Debug("generating default encryption configuration",
	)

	return e, nil
}

