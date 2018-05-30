package data

import (
	"io/ioutil"
	"encoding/json"
	"errors"
	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/source"
	"os"
)

var countries model.Countries

func init() {
	jsonData, err := ioutil.ReadFile("data/sources/sources.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonData, &countries)
	if err != nil {
		panic(err)
	}
}

func GetCountries() []string {
	var names []string

	for _, country := range countries {
		names = append(names, country.Name)
	}

	return names
}

func getModelFile(countryCode string) (string, error) {
	for _, country := range countries {
		if country.Name == countryCode {
			return country.File, nil
		}
	}

	return "", errors.New("country file not found")
}

func GetModel(countryCode string) (source.Model, error) {
	var countryFile, localFile string
	var countryModel source.Model

	countryFile, err := getModelFile(countryCode)
	if err != nil {
		return countryModel, err
	}

	localFile = "data/sources/" + countryFile
	if _, err := os.Stat(localFile); os.IsNotExist(err) {
		return countryModel, os.ErrNotExist
	}

	jsonData, err := ioutil.ReadFile(localFile)
	if err != nil {
		return countryModel, err
	}

	err = json.Unmarshal(jsonData, &countryModel)
	if err != nil {
		return countryModel, err
	}

	return countryModel, nil
}

func GetAttributes( countryCode string ) []string {
	var names []string
	var countryModel source.Model

	countryModel, err := GetModel( countryCode )
	if err != nil {
		panic(err)
	}

	for _, attribute := range countryModel.Attributes {
		names = append(names, attribute.Name)
	}

	return names
}
