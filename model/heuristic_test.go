package model

import (
	"testing"
	"fmt"
	"strings"
)

func TestValidateHeuristic(t *testing.T) {
	var valid = [1]string{"63b347e7-512f-434c-aae3-39e6e9455899"}

	for _, item := range valid {
		if ok, err := ValidateHeuristic( strings.ToLower( item ) ); !ok || err != nil {
			s := fmt.Sprintf("expected uuid v4 formatted heuristic ('%s')  = ok", item )
			t.Error( s )
		}

		if ok, err := ValidateHeuristic( strings.ToUpper( item ) ); !ok || err != nil {
			s := fmt.Sprintf("expected UUID v4 formatted heuristic ('%s')  = ok", item )
			t.Error( s )
		}
	}

	var invalid = [2]string{"63b347e7512f434caae339e6e9455899", "9cbe45d0-7e0d-11e8-adc0-fa7ae01bbebc"}

	for _, item := range invalid {
		if ok, err := ValidateHeuristic( strings.ToLower( item ) ); ok || err == nil {
			s := fmt.Sprintf("expected uuid v4 formatted heuristic ('%s') != ok", item )
			t.Error( s )
		}

		if ok, err := ValidateHeuristic( strings.ToUpper( item ) ); ok || err == nil {
			s := fmt.Sprintf("expected UUID v4 formatted heuristic ('%s') != ok", item )
			t.Error( s )
		}
	}
}
