package server

import (
	"net/http"
	"time"
	"log"
	"encoding/json"
	"fmt"
)

type LogEntry struct {
	Time time.Time `json:"time"`
	Method string `json:"method"`
	URI string `json:"uri"`
	Name string `json:"name"`
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc( func( w http.ResponseWriter, r *http.Request){
		var logEntry LogEntry
		logEntry = LogEntry{Time: time.Now().UTC(), Method: r.Method, URI: r.RequestURI, Name: name}

		entryJson, err := json.Marshal(logEntry)
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf( "%s", string(entryJson))

		inner.ServeHTTP(w,r)
	})
}
