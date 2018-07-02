package model

import (
	"testing"
	"fmt"
)

func TestValidateCountryCode(t *testing.T) {
	var valid = [2]string{"us", "US"}

	for _, item := range valid {
		if ok, err := ValidateCountryCode( item ); !ok || err != nil {
			s := fmt.Sprintf("expected country code ('%s')  = ok", item )
			t.Error( s )
		}


	}

	var invalid = [3]string{"", "u", "usa"}

	for _, item := range invalid {
		if ok, err := ValidateCountryCode( item ); ok || err == nil {
			s := fmt.Sprintf("expected country code ('%s') != ok", item )
			t.Error( s )
		}
	}
}
