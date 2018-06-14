package sys

import (
	"net/http"
	"encoding/json"
	"time"
	"strings"

	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/source"
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
	r = append( r, model.Route{ Name: "SysHealth", Method: "GET", Pattern: "/v1/sys/health", HandlerFunc: Health} )
	r = append( r, model.Route{ Name: "SysReload", Method: "GET", Pattern: "/v1/sys/reload", HandlerFunc: Reload} )

	return r
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader( http.StatusOK )

	var SysInfo SysHealth
	SysInfo.ApiVersion = VERSION

	var m *source.Model
	var err error
	var countries []string
	countries = data.GetCountries()
	for _, country := range countries {
		m, err = data.GetModel(country)
		if err == nil {
			m := ModelVersion{ CountryCode: strings.ToUpper( country ), Timestamp: m.ModelDate.UTC(), Version: m.ModelVersion}
			SysInfo.ModelVersions = append( SysInfo.ModelVersions, m)
		}
	}

	if err := json.NewEncoder(w).Encode(SysInfo); err != nil {
		panic(err)
	}
}

func Reload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")

	if err := data.Reload(); err == nil {
		w.WriteHeader( http.StatusOK )
	} else {
		w.WriteHeader( http.StatusInternalServerError )
	}
}
