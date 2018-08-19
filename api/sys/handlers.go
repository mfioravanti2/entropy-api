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
	"github.com/mfioravanti2/entropy-api/data/scoringdb"
	"github.com/mfioravanti2/entropy-api/model/metrics"
	"github.com/mfioravanti2/entropy-api/config"
	"github.com/mfioravanti2/entropy-api/command/server/enforce"
)

type DataStore struct {
	Status		string		`json:"status"`
	Engine		string		`json:"engine"`
	LastUse		time.Time	`json:"last_use"`
}

type ModelVersion struct {
	CountryCode string		`json:"country"`
	Timestamp   time.Time	`json:"timestamp"`
	Version     string		`json:"version"`
}

type ModelVersions []ModelVersion

type SysHealth struct {
	Status			string			`json:"status"`
	ApiVersion		string			`json:"api_version"`
	ModelVersions	ModelVersions	`json:"model_versions"`
	DataStore		DataStore		`json:"data_store"`
}

const (
	VERSION = "0.0.1"
)

// Add Handlers for the System Configuration/Health Endpoints
func AddHandlers(r model.Routes, endpoints *config.Endpoints) model.Routes {
	ctx := logging.WithFuncId( context.Background(), "AddHandlers", "sys" )

	logger := logging.Logger( ctx )

	endpoint, err := endpoints.GetEndpoint( config.ENDPOINT_HEALTH )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", config.ENDPOINT_HEALTH ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers",
				zap.String( "endpoint", "/v1/sys/health" ),
				)
			r = append( r, model.Route{ Name: "SysHealth",
										Method: "GET",
										Pattern: "/v1/sys/health",
										HandlerFunc: Health,
										Params: nil,
										Enforce: model.ENFORCE_CONTENT_NONE,
										Policy: endpoint,
										AuthN: model.AUTH_METHOD_NONE } )
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/sys/health" ),
				zap.String( "policy", config.ENDPOINT_HEALTH ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", config.ENDPOINT_HEALTH ),
		)
	}

	endpoint, err = endpoints.GetEndpoint( config.ENDPOINT_SYSTEM )
	if err == nil {
		logger.Info("checking handler endpoint policy",
			zap.String( "policy", config.ENDPOINT_SYSTEM ),
			zap.Bool( "enabled", endpoint.Enabled ),
		)

		if endpoint.Enabled {
			logger.Debug("registering handlers",
				zap.String( "endpoint", "/v1/sys/reload" ),
				)
			r = append( r, model.Route{ Name: "SysReload",
										Method: "GET",
										Pattern: "/v1/sys/reload",
										HandlerFunc: Reload,
										Params: nil,
										Enforce: model.ENFORCE_CONTENT_NONE,
										Policy: endpoint,
										AuthN: model.AUTH_METHOD_NONE } )
		} else {
			logger.Warn("handler disabled by configuration",
				zap.String( "endpoint", "/v1/sys/reload" ),
				zap.String( "policy", config.ENDPOINT_SYSTEM ),
			)
		}
	} else {
		logger.Error("unable to locate endpoint policy",
			zap.String( "policy", config.ENDPOINT_SYSTEM ),
		)
	}

	return r
}

// Return the system's health
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", enforce.HEADER_JSON_CONTENT_TYPE)
	w.WriteHeader( http.StatusOK )

	var SysInfo SysHealth
	SysInfo.ApiVersion = VERSION

	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.sys.health.get" )
	ctrReg.Inc(1)

	logger.Info( "checking system health",
		zap.String("version", SysInfo.ApiVersion ),
	)

	var errCount = 0
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
			errCount += 1

			logger.Error( "checking model health",
				zap.String("model ", strings.ToUpper( country ) ),
				zap.Bool( "loaded", false ),
				zap.String( "status", "error" ),
				zap.String("error ", err.Error() ),
			)
		}
	}

	ds, err := scoringdb.GetDataStore( nil )
	if ds != nil {
		// retrieve the data store engine
		SysInfo.DataStore.Engine = strings.ToLower( ds.Config().Engine )
		SysInfo.DataStore.LastUse = ds.LastUse

		// ping the data store to see if it is available
		if ds.Ready( false ) {
			SysInfo.DataStore.Status = "good"
		} else {
			SysInfo.DataStore.Status = "degraded"
		}
	} else {
		SysInfo.DataStore.Status = "offline"
		SysInfo.DataStore.Engine = "unknown"
		SysInfo.DataStore.LastUse = time.Now().UTC()
	}

	// if no errors occurred and the data store is operational, the endpoint is good
	SysInfo.Status = "degraded"
	if errCount == 0 {
		if SysInfo.DataStore.Status == "good" {
			SysInfo.Status = "good"
		} else {
			SysInfo.Status = "degraded"
		}
	}

	// encode and return the response to the client
	if err := json.NewEncoder(w).Encode(SysInfo); err != nil {
		logger.Error( "encoding system health",
			zap.String( "status", "error" ),
			zap.String("error", err.Error() ),
		)
	}

	ctrReg, _ = metrix.GetCounter( "entropy.sys.health.get.status.200" )
	ctrReg.Inc(1)
}

// Reload the country models
func Reload(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()
	logger := logging.Logger(reqCtx)

	ctrReg, _ := metrix.GetCounter( "entropy.sys.reload.get" )
	ctrReg.Inc(1)

	// Reload the country models
	if err := data.Reload( reqCtx ); err == nil {
		w.Header().Set("Content-Type", enforce.HEADER_JSON_CONTENT_TYPE)
		w.WriteHeader( http.StatusOK )

		logger.Info( "reloading models",
			zap.String( "status", "ok" ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.sys.reload.get.status.200" )
		ctrReg.Inc(1)
	} else {
		w.WriteHeader( http.StatusInternalServerError )

		logger.Error( "reloading models",
			zap.String( "status", "error" ),
			zap.String( "error", err.Error() ),
		)

		ctrReg, _ := metrix.GetCounter( "entropy.sys.reload.get.status.500" )
		ctrReg.Inc(1)
	}
}
