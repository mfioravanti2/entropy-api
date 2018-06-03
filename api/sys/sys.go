package sys

import (
	"net/http"
	"encoding/json"
	"github.com/mfioravanti2/entropy-api/model"
)

type SysHealth struct {
	ApiVersion string `json:"api_version"`
}

var SysInfo = SysHealth{"0.0.1"}

func AddHandlers(r model.Routes) model.Routes {
	r = append( r, model.Route{"SysHealth", "GET", "/v1/sys/health", Health} )

	return r
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader( http.StatusOK )

	if err := json.NewEncoder(w).Encode(SysInfo); err != nil {
		panic(err)
	}
}


