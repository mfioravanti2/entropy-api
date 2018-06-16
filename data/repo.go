package data

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"errors"
	"os"
	"fmt"
	"sort"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/source"
)

var modelCache map[string]*source.Model

func init() {
	if err := Reload(); err != nil {
		panic(err)
	}
}

func Reload() error {
	jsonData, err := ioutil.ReadFile("data/sources/sources.json")
	if err != nil {
		s := fmt.Sprintf("unable to load source file")
		return errors.New(s)
	}

	var countries model.Countries
	err = json.Unmarshal(jsonData, &countries)
	if err != nil {
		s := fmt.Sprintf("unable to parse source file, expected json format")
		return errors.New(s)
	}

	modelCache = make( map[string]*source.Model )
	for _, country := range countries {
		countryCode := strings.ToLower( country.Name )
		localFile := "data/sources/" + country.File
		if _, err := os.Stat(localFile); os.IsNotExist(err) {
			s := fmt.Sprintf("country model file (%s), does not exist", country.File)
			return errors.New(s)
		}

		jsonData, err := ioutil.ReadFile(localFile)
		if err != nil {
			s := fmt.Sprintf("unable to read country model file (%s)", country.File)
			return errors.New(s)
		}

		var countryModel source.Model
		err = json.Unmarshal(jsonData, &countryModel)
		if err != nil {
			s := fmt.Sprintf("unable to parse country model file (%s), expected json format", country.File)
			return errors.New(s)
		}

		modelCache[countryCode] = &countryModel
	}

	return nil
}

func GetCountries() []string {
	var names []string

	for countryCode, _ := range modelCache {
		names = append(names, countryCode)
	}

	return names
}

func GetModel(countryCode string) (*source.Model, error) {
	var countryModel *source.Model
	if len(countryCode) != 2 {
		s := fmt.Sprintf("country code not specified")
		return countryModel, errors.New(s)
	}

	if countryModel, ok := modelCache[countryCode]; ok {
		return countryModel, nil
	}

	s := fmt.Sprintf("country model (%s) not found", strings.ToUpper(countryCode))
	return countryModel, errors.New(s)
}

func GetAttributes( countryCode string ) ([]string, error) {
	var names []string
	var countryModel *source.Model

	countryModel, err := GetModel( countryCode )
	if err != nil {
		return names, err
	}

	for _, attribute := range countryModel.Attributes {
		names = append(names, attribute.Mnemonic)
	}

	sort.Strings(names)

	return names, nil
}

func GetAttribute( countryCode string, attributeMnemonic string ) (source.Attribute, error) {
	var err error
	var attribute source.Attribute

	countryModel, err := GetModel( countryCode )
	if err != nil {
		return attribute, err
	}

	for _, attribute := range countryModel.Attributes {
		if attribute.Mnemonic == attributeMnemonic {
			return attribute, nil
		}
	}

	s := fmt.Sprintf("attribute (%s) for country (%s) not found", attributeMnemonic, strings.ToUpper(countryCode))
	return attribute, errors.New(s)
}
