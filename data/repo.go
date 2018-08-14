package data

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"errors"
	"os"
	"fmt"
	"sort"
	"context"

	"go.uber.org/zap"

	"github.com/mfioravanti2/entropy-api/model"
	"github.com/mfioravanti2/entropy-api/model/source"
	"github.com/mfioravanti2/entropy-api/command/server/logging"
)

var modelCache map[string]*source.Model

func init() {
	ctx := logging.WithFuncId( context.Background(), "init", "data" )

	logger := logging.Logger( ctx )

	logger.Debug("loading models" )

	if err := Reload( ctx ); err != nil {
		logger.Error( "loading models",
			zap.String("error", err.Error() ),
		)
	}
}

func Reload( ctx context.Context ) error {
	if ctx == nil {
		ctx = logging.WithFuncId( context.Background(), "Reload", "data" )
	} else {
		ctx = logging.WithFuncId( ctx, "Reload", "data" )
	}

	logger := logging.Logger( ctx )

	logger.Debug("preparing to reloading models",
		zap.String("model_file", "data/sources/sources.json" ),
	)

	jsonData, err := ioutil.ReadFile("data/sources/sources.json")
	if err != nil {
		s := fmt.Sprintf("unable to load source file")
		logger.Error( "loading model index",
			zap.String("file", "data/sources/sources.json" ),
			zap.String("error", s ),
		)

		return errors.New(s)
	}

	logger.Debug("attempting to unmarshal model index",
		zap.String("model_file", "data/sources/sources.json" ),
	)

	var countries model.Countries
	err = json.Unmarshal(jsonData, &countries)
	if err != nil {
		s := fmt.Sprintf("unable to parse source file, expected json format")
		logger.Error( "unable to parse model index",
			zap.String("file", "data/sources/sources.json" ),
			zap.String("error", s ),
		)

		return errors.New(s)
	}

	modelCache = make( map[string]*source.Model )
	for _, country := range countries {
		countryCode := strings.ToLower( country.Name )
		localFile := "data/sources/" + country.File

		logger.Debug( "loading country model",
			zap.String("file", localFile ),
			zap.String("countryId", strings.ToUpper( countryCode ) ),
		)

		if _, err := os.Stat(localFile); os.IsNotExist(err) {
			s := fmt.Sprintf("country model file (%s), does not exist", country.File)
			logger.Error( "unable to locate country model",
				zap.String("file", localFile ),
				zap.String("countryId", strings.ToUpper( countryCode ) ),
				zap.String("error", s ),
			)

			return errors.New(s)
		}

		jsonData, err := ioutil.ReadFile(localFile)
		if err != nil {
			s := fmt.Sprintf("unable to read country model file (%s)", country.File)
			logger.Error( "unable to read country model file",
				zap.String("file", localFile ),
				zap.String("countryId", strings.ToUpper( countryCode ) ),
				zap.String("error", s ),
			)

			return errors.New(s)
		}

		var countryModel source.Model
		err = json.Unmarshal(jsonData, &countryModel)
		if err != nil {
			s := fmt.Sprintf("unable to parse country model file (%s), expected json format", country.File)
			logger.Error( "unable to unmarshal country model file",
				zap.String("file", localFile ),
				zap.String("countryId", strings.ToUpper( countryCode ) ),
				zap.String("error", s ),
			)

			return errors.New(s)
		}

		logger.Debug( "registering country model",
			zap.String("file", localFile ),
			zap.String("countryId", strings.ToUpper( countryCode ) ),
		)

		modelCache[countryCode] = &countryModel
	}

	logger.Info( "country models loaded" )

	return nil
}

func GetCountries() []string {
	var names []string

	for countryCode, _ := range modelCache {
		names = append(names, countryCode)
	}

	sort.Strings(names)

	return names
}

func GetAllCountries() []source.Model {
	var countries []source.Model

	for countryCode, _ := range modelCache {
		var country *source.Model
		country = modelCache[countryCode]

		countries = append(countries, *country)
	}

	return countries
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

	names, err = countryModel.Attributes.ToStrings()
	return names, err
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

func GetHeuristics( countryCode string ) ([]string, error) {
	var names []string
	var countryModel *source.Model

	countryModel, err := GetModel( countryCode )
	if err != nil {
		return names, err
	}

	names, err = countryModel.Heuristics.ToStrings()
	return names, nil
}

func GetHeuristic( countryCode string, heuristicId string ) (source.Heuristic, error) {
	var err error
	var heuristic source.Heuristic

	countryModel, err := GetModel( countryCode )
	if err != nil {
		return heuristic, err
	}

	for _, heuristic := range countryModel.Heuristics {
		if heuristic.Id == heuristicId {
			return heuristic, nil
		}
	}

	s := fmt.Sprintf("heuristic (%s) for country (%s) not found", heuristicId, strings.ToUpper(countryCode))
	return heuristic, errors.New(s)
}
