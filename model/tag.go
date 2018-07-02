package model

import (
	"regexp"
	"strings"
	"fmt"
	"errors"
)

func ValidateTag( tag string ) (bool, error) {
	var err error

	if tag == "" {
		return true, nil
	}

	rx, err := regexp.Compile( `^[A-Z0-9]+$` )
	if err != nil {
		return false, err
	}

	if rx.MatchString( strings.ToUpper(tag) ) {
		return true, nil
	} else {
		s := fmt.Sprintf("tag (%s) failed validation", tag )
		return false, errors.New(s)
	}

	return false, nil
}


