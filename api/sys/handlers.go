package sys

import (
	"context"
	"net/http"
	"encoding/json"
	"time"
	"strings"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/source"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

type ModelVersion struct {
	CountryCode string    `json:"country"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
}

type ModelVersions []ModelVersion

type SysHealth struct {
	ApiVersion    string        `json:"api_version"`
	ModelVersions ModelVersions `json:"model_versions"`
}

const (
	VERSION = "0.0.1"
)

func AddHandlers(r model.Routes) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "sys" )

	logger := logging.Logger( ctx )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/sys/health" ) )
	r = append( r, model.Route{ Name: "SysHealth", Method: "GET", Pattern: "/v1/sys/health", HandlerFunc: Health} )

	logger.Debug("registering handlers", zap.String( "endpoint", "/v1/sys/reload" ) )
	r = append( r, model.Route{ Name: "SysReload", Method: "GET", Pattern: "/v1/sys/reload", HandlerFunc: Reload} )

	return r
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader( http.StatusOK )

	var SysInfo SysHealth
	SysInfo.ApiVersion = VERSION

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	logger.Info( "checking system health",
		zap.String("version", SysInfo.ApiVersion ),
	)

	var m *source.Model
	var err error
	var countries []string
	countries = data.GetCountries()
	for _, country := range countries {
		m, err = data.GetModel(country)
		if err == nil {
			m := ModelVersion{ CountryCode: strings.ToUpper( country ), Timestamp: m.ModelDate.UTC(), Version: m.ModelVersion}
			SysInfo.ModelVersions = append( SysInfo.ModelVersions, m)

			logger.Info( "checking model health",
				zap.String("model ", strings.ToUpper( country ) ),
				zap.String( "status", "ok" ),
				zap.Bool( "loaded", true ),
			)
		} else {
			logger.Error( "checking model health",
				zap.String("model ", strings.ToUpper( country ) ),
				zap.Bool( "loaded", false ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}
	}

	if err := json.NewEncoder(w).Encode(SysInfo); err != nil {
		logger.Error( "encoding system health",
			zap.String( "status", "error" ),
			zap.String("error", err.Error() ),
		)
	}
}

func Reload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	if err := data.Reload( reqCtx ); err == nil {
		w.WriteHeader( http.StatusOK )

		logger.Info( "reloading models",
			zap.String( "status", "ok" ),
		)
	} else {
		w.WriteHeader( http.StatusInternalServerError )

		logger.Error( "reloading models",
			zap.String( "status", "error" ),
			zap.String( "error", err.Error() ),
		)
	}
}