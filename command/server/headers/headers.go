package headers

import (
	"fmt"
	"errors"
	"net/http"
	"context"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type Header struct {
	Name string
	Value string
}

type Headers []Header
var headers Headers

var hCacheCtrl  = Header{Name: "Cache-Control", Value: "nocache, nostore, mustrevalidate"}
var hCacheExpr  = Header{Name: "Expires", Value: "0"}
var hXSSProtect = Header{Name: "X-XSS-Protection", Value: "1; mode=block"}
var hXContent   = Header{Name: "X-Content-Type-Options", Value: "nosniff"}
var hXFramOpts  = Header{Name: "X-Frame-Options", Value: "SAMEORIGIN"}

var hCORSOrgin  = Header{Name: "Access-Control-Allow-Origin", Value: "*"}


func init() {
	headers = Headers{}
}

func BuildHeaders( cfgHeaders *config.Headers ) error {
	ctx := logging.WithFuncId(context.Background(), "BuildHeaders", "headers" )

	logger := logging.Logger(ctx)
	logger.Info("building security headers",
	)

	var count = 0
	headers = Headers{}

	for _, cfgHeader := range *cfgHeaders {
		logger.Debug("validating security header",
			zap.String("Name", cfgHeader.Name),
			zap.Bool("Enabled", cfgHeader.Enabled ),
			zap.String("Field", cfgHeader.Field),
			zap.String("Value", cfgHeader.Value ),
		)

		if cfgHeader.Enabled {
			if cfgHeader.Value != "" {
				logger.Debug("registering security header",
					zap.String("Name", cfgHeader.Name),
					zap.String("Field", cfgHeader.Field),
					zap.String("Value", cfgHeader.Value ),
				)

				headers = append( headers, Header{ Name: cfgHeader.Field, Value: cfgHeader.Value } )
				count++
			} else {
				logger.Error("security header found with no value",
					zap.String("Name", cfgHeader.Name),
					zap.String("Field", cfgHeader.Field),
					zap.String("Value", cfgHeader.Value ),
				)

				s := fmt.Sprintf("security hearer with no value (%s: %s)", cfgHeader.Name, cfgHeader.Field )
				return errors.New(s)
			}
		} else {
			logger.Debug("disabled security header skipped",
				zap.String("Name", cfgHeader.Name),
				zap.String("Field", cfgHeader.Field),
				zap.String("Value", cfgHeader.Value ),
			)
		}
	}

	if count == 0 {
		logger.Warn("completed registering security headers",
			zap.Int( "count", count ),
		)
	}

	return nil
}

func SecurityHeadersHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range headers {
			w.Header().Set( h.Name, h.Value )
		}

		next.ServeHTTP(w, r)
	})
}
