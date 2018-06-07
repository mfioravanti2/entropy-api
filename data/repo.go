package data

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"errors"
	"os"
	"fmt"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/source"
)

var modelCache map[string]*source.Model

func init() {
	reload()
}

func reload() {
	jsonData, err := ioutil.ReadFile("data/sources/sources.json")
	if err != nil {
		panic(err)
	}

	var countries model.Countries
	err = json.Unmarshal(jsonData, &countries)
	if err != nil {
		panic(err)
	}

	modelCache = make( map[string]*source.Model )
	for _, country := range countries {
		countryCode := strings.ToLower( country.Name )
		localFile := "data/sources/" + country.File
		if _, err := os.Stat(localFile); os.IsNotExist(err) {
			panic(os.ErrNotExist)
		}

		jsonData, err := ioutil.ReadFile(localFile)
		if err != nil {
			panic(err)
		}

		var countryModel source.Model
		err = json.Unmarshal(jsonData, &countryModel)
		if err != nil {
			panic(err)
		}

		modelCache[countryCode] = &countryModel
	}
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

func GetAttributes( countryCode string ) []string {
	var names []string
	var countryModel *source.Model

	countryModel, err := GetModel( countryCode )
	if err != nil {
		panic(err)
	}

	for _, attribute := range countryModel.Attributes {
		names = append(names, attribute.Name)
	}

	return names
}
