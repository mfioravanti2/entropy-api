package model

import (
	"testing"
	"fmt"
)

func TestValidateHeuristic(t *testing.T) {
	var valid = [2]string{"63b347e7-512f-434c-aae3-39e6e9455899", "63B347E7-512F-434C-AAE3-39E6E9455899"}

	for _, item := range valid {
		if ok, err := ValidateHeuristic( item ); !ok || err != nil {
			s := fmt.Sprintf("expected uuid v4 formatted heuristic ('%s')  = ok", item )
			t.Error( s )
		}
	}

	var invalid = [2]string{"63b347e7512f434caae339e6e9455899", "63B347E7512F434CAAE339E6E9455899"}

	for _, item := range invalid {
		if ok, err := ValidateHeuristic( item ); ok || err == nil {
			s := fmt.Sprintf("expected uuid v4 formatted heuristic ('%s') != ok", item )
			t.Error( s )
		}
	}
}
