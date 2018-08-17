package cli

import (
	"context"
	"fmt"
	"errors"

	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

const (
	HEADER_CACHE_CTRL	string = "entropy.security.headers.cache_control"
	HEADER_EXPIRES		string = "entropy.security.headers.expires"
	HEADER_XSS_PROTECT	string = "entropy.security.headers.xss_protect"
	HEADER_X_CONTENT	string = "entropy.security.headers.x_content_type"
	HEADER_X_FRM_OPTS	string = "entropy.security.headers.x_frame_opts"
	HEADER_CORS_ORGIN	string = "entropy.security.headers.cors.origin"
	HEADER_C_SEC_POLY	string = "entropy.security.headers.policy.content"
	HEADER_RFERR_POLY	string = "entropy.security.headers.policy.referrer"
	HEADER_FEATR_POLY	string = "entropy.security.headers.policy.feature"
)

type Header struct {
	Name string		`json:"name"`
	Enabled bool	`json:"enabled"`
	Field string	`json:"field"`
	Value string	`json:"value"`
}

type Headers []Header

func (h *Header) DefaultConfig() {
	h.Name = ""
	h.Enabled = true
	h.Field = ""
	h.Value = ""
}

// Create a default empty path
func NewHeader() (*Header, error) {
	var h = &Header{}
	h.DefaultConfig()

	ctx := logging.WithFuncId( context.Background(), "NewHeader", "cli" )

	logger := logging.Logger( ctx )
	logger.Debug("generating default header configuration",
	)

	return h, nil
}

func NewHeaders() (Headers, error) {
	var hs = Headers{}
	var h *Header

	h, _ = NewHeader()
	h.Name = HEADER_CACHE_CTRL
	h.Field = "Cache-Control"
	h.Value = "nocache, nostore, mustrevalidate"
	h.Enabled = true
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_EXPIRES
	h.Field = "Expires"
	h.Value = "0"
	h.Enabled = true
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_XSS_PROTECT
	h.Field = "X-XSS-Protection"
	h.Value = "1; mode=block"
	h.Enabled = true
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_X_CONTENT
	h.Field = "X-Content-Type-Options"
	h.Value = "nosniff"
	h.Enabled = true
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_X_FRM_OPTS
	h.Field = "X-Frame-Options"
	h.Value = "SAMEORIGIN"
	h.Enabled = true
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_CORS_ORGIN
	h.Field = "Access-Control-Allow-Origin"
	h.Value = "*"
	h.Enabled = true
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_C_SEC_POLY
	h.Field = "Content-Security-Policy"
	h.Value = "default-src 'none'; script-src 'self'; img-src 'self'; style-src 'self'"
	h.Enabled = true
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_RFERR_POLY
	h.Field = "Referrer-Policy"
	h.Value = "strict-origin"
	h.Enabled = false
	hs = append( hs, *h )

	h, _ = NewHeader()
	h.Name = HEADER_FEATR_POLY
	h.Field = "Feature-Policy"
	h.Value = "sync-xhr 'self'"
	h.Enabled = false
	hs = append( hs, *h )

	return hs, nil
}

func (ps *Paths) GetHeader( name string ) (*Path, error) {
	var p Path

	for _, p = range *ps {
		if p.Name == name {
			return &p, nil
		}
	}

	s := fmt.Sprintf("find header failed (header name: %s)", name )
	return nil, errors.New(s)
}
