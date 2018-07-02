package model

import (
	"testing"
	"fmt"
)

func TestValidateTag(t *testing.T) {
	var valid = [3]string{"work", "personal", ""}

	for _, item := range valid {
		if ok, err := ValidateTag( item ); !ok || err != nil {
			s := fmt.Sprintf("expected tag ('%s')  = ok", item )
			t.Error( s )
		}
	}

	var invalid = [3]string{"work+", "personal.", "-"}

	for _, item := range invalid {
		if ok, err := ValidateTag( item ); ok || err == nil {
			s := fmt.Sprintf("expected tag ('%s') != ok", item )
			t.Error( s )
		}
	}
}
