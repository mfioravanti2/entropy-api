package model

import (
	"regexp"
	"fmt"
	"errors"
)

type Country struct {
	Name string `json:"name"`
	File string `json:"file"`
}

type Countries []Country

func ValidateCountryCode( countryCode string ) (bool, error) {
	var err error

	// two-digit country codes: ISO 3166-1 alpha-2
	rx, err := regexp.Compile(`^[a-zA-Z]{2}$` )
	if err != nil {
		return false, err
	}

	if rx.MatchString( countryCode ) {
		return true, nil
	} else {
		s := fmt.Sprintf("country code (%s) failed validation", countryCode )
		return false, errors.New(s)
	}

	return false, nil
}
