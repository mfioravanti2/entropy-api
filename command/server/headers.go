package server

import "net/http"

type Header struct {
	Name string
	Value string
}

var hCacheCtrl  = Header{Name: "Cache-Control", Value: "nocache, nostore, mustrevalidate"}
var hCacheExpr  = Header{Name: "Expires", Value: "0"}
var hXSSProtect = Header{Name: "X-XSS-Protection", Value: "1; mode=block"}
var hXContent   = Header{Name: "X-Content-Type-Options", Value: "nosniff"}
var hXFramOpts  = Header{Name: "X-Frame-Options", Value: "SAMEORIGIN"}

func SecurityHeaders( w http.ResponseWriter ) {

	// Setup HTTP Response Cache Control options
	w.Header().Set( hCacheCtrl.Name, hCacheCtrl.Value )
	w.Header().Set( hCacheExpr.Name, hCacheCtrl.Value )

	// Setup XSS Protection options
	w.Header().Set( hXSSProtect.Name, hXSSProtect.Value )

	// Setup X Content Protection options
	w.Header().Set( hXContent.Name, hXContent.Value )

	// Setup X-Frame options
	w.Header().Set( hXFramOpts.Name, hXFramOpts.Value )
}
