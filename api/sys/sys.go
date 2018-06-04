package sys

import (
	"net/http"
	"encoding/json"
	"github.com/mfioravanti2/entropy-api/model"
	"time"
	"github.com/mfioravanti2/entropy-api/data"
	"github.com/mfioravanti2/entropy-api/model/source"
)

type ModelVersion struct {
	CountryCode string `json:"country"`
	Timestamp time.Time `json:"timestamp"`
	Version string `json:"version"`
}

type ModelVersions []ModelVersion

type SysHealth struct {
	ApiVersion string `json:"api_version"`
	ModelVersions ModelVersions `json:"model_versions"`
}

const (
	VERSION = "0.0.1"
)

//var SysInfo = SysHealth{"0.0.1"}

func AddHandlers(r model.Routes) model.Routes {
	r = append( r, model.Route{"SysHealth", "GET", "/v1/sys/health", Health} )

	return r
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader( http.StatusOK )

	var SysInfo SysHealth
	SysInfo.ApiVersion = VERSION

	var m source.Model
	var err error
	var countries []string
	countries = data.GetCountries()
	for _, country := range countries {
		m, err = data.GetModel(country)
		if err == nil {
			SysInfo.ModelVersions = append( SysInfo.ModelVersions, ModelVersion{ country, m.ModelDate, m.ModelVersion})
		}
	}

	if err := json.NewEncoder(w).Encode(SysInfo); err != nil {
		panic(err)
	}
}


