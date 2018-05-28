package data

import (
	"io/ioutil"
	"github.com/mfioravanti2/entropy-api/model"
	"encoding/json"
)

// var countries map[string]string
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

func GetAttributes( country string ) model.Attributes {
	var attributes model.Attributes

	return attributes
}
