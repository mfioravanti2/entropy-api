package model

import (
	"regexp"
	"strings"
	"fmt"
	"errors"
)

func ValidateFormat( formatId string ) (bool, error) {
	var err error

	rx, err := regexp.Compile( `^(mean|naive|rare)$` )
	if err != nil {
		return false, err
	}

	if rx.MatchString( strings.ToLower(formatId) ) {
		return true, nil
	} else {
		s := fmt.Sprintf("formatId (%s) failed validation", formatId )
		return false, errors.New(s)
	}

	return false, nil
}

