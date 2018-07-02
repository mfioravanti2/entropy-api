package model

import (
	"fmt"
	"testing"
)

func TestValidateFormat(t *testing.T) {
	var list = [3]string{"mean", "naive", "rare"}

	for _, item := range list {
		if ok, err := ValidateFormat( item ); !ok || err != nil {
			s := fmt.Sprintf("expected format code(%s) = ok", item)
			t.Error( s )
		}
	}
}
