package sys

import (
	"net/http"
	"encoding/json"
	"github.com/mfioravanti2/entropy-api/model"
)

type SysHealth struct {
	Version string `json:"version"`
}

func AddHandlers(r model.Routes) model.Routes {
	r = append( r, model.Route{"SysHealth", "GET", "/v1/sys/health", Health} )

	return r
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader( http.StatusOK )

	if err := json.NewEncoder(w).Encode(SysHealth{"0.0.1"}); err != nil {
		panic(err)
	}
}


