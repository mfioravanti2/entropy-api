package model

import (
	"regexp"
	"strings"
	"fmt"
	"errors"
)

func ValidateHeuristic( heuristicId string ) (bool, error) {
	var err error

	rx, err := regexp.Compile( `^[0-9A-F]{8}-[0-9A-F]{4}-4[0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$` )
	if err != nil {
		return false, err
	}

	if rx.MatchString( strings.ToUpper(heuristicId) ) {
		return true, nil
	} else {
		s := fmt.Sprintf("heurisitic (%s) failed validation", heuristicId )
		return false, errors.New(s)
	}

	return false, nil
}

