package data

import (
	"io/ioutil"
	"github.com/mfioravanti2/entropy-api/model"
	"encoding/json"
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

func GetAttributes( countryCode string ) []string {
	var names []string

	for _, country := range countries {
		if country.Name == countryCode {
			jsonData, err := ioutil.ReadFile("data/sources/" + country.File)
			if err != nil {
				panic(err)
			}

			var entropy model.Entropy

			err = json.Unmarshal(jsonData, &entropy)
			if err != nil {
				panic(err)
			}

			for _, attribute := range entropy.Attributes {
				names = append(names, attribute.Name)
			}

			break
		}
	}


	return names
}
