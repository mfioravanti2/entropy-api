package model

import (
	"regexp"
	"fmt"
	"errors"
)

type Attribute struct {
	Name string `json:"name"`
}

type Attributes []Attribute

func ValidateAttributeMnemonic( attributeId string ) (bool, error) {
	var err error

	rx, err := regexp.Compile(`^([a-zA-Z0-9_]+.)+([a-zA-Z0-9_])$` )
	if err != nil {
		return false, err
	}

	if rx.MatchString( attributeId ) {
		return true, nil
	} else {
		s := fmt.Sprintf("attribute.mnemnoic (%s) failed validation", attributeId )
		return false, errors.New(s)
	}

	return false, err
}


