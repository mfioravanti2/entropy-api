package model

import (
	"testing"
	"fmt"
)

func TestValidateAttributeMnemonic(t *testing.T) {
	var valid = [2]string{"date_of_birth", "postal_address.plus_4_digit_code"}

	for _, item := range valid {
		if ok, err := ValidateAttributeMnemonic( item ); !ok || err != nil {
			s := fmt.Sprintf("expected attribute mnemonic ('%s')  = ok", item )
			t.Error( s )
		}
	}

	var invalid = [3]string{"postal_address.plus_4_digit_code.", "postal_address.+4_digit_code", ""}

	for _, item := range invalid {
		if ok, err := ValidateAttributeMnemonic( item ); ok || err == nil {
			s := fmt.Sprintf("expected attribute mnemonic ('%s') != ok", item )
			t.Error( s )
		}
	}
}
