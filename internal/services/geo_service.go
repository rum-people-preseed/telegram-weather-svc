package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/biter777/countries"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services/utils"
)

type GeoService interface {
	ValidateCountry(country string) (string, error)
	ValidateCity(city string, country string) (string, models.Coordinates, error)
}

type GeoNameService struct {
	preparedURL string
	sender      Sender
	log         models.Logger
	fuzzy       string
}

func NewGeoNameService(sender Sender, logger models.Logger, fuzzy float64) GeoNameService {
	baseURL := "http://api.geonames.org/searchJSON"
	apiKey := os.Getenv("GEO_NAME_SERVICE_USERNAME")
	apiKeyParam := utils.NewHTTPParam("username", apiKey)
	maxRowsParam := utils.NewHTTPParam("maxRows", "1")
	preparedURL := utils.BuildURL(baseURL, apiKeyParam, maxRowsParam)
	fuzzyStr := strconv.FormatFloat(fuzzy, 'g', 2, 64)
	return GeoNameService{sender: sender, preparedURL: preparedURL, log: logger, fuzzy: fuzzyStr}
}

func (s *GeoNameService) ValidateCountry(country string) (string, error) {
	featureClassParam := utils.NewHTTPParam("featureClass", "A")
	nameEqualsParam := utils.NewHTTPParam("name_equals", country)
	fuzzyParam := utils.NewHTTPParam("fuzzy", s.fuzzy)

	URL := utils.BuildURL(s.preparedURL, featureClassParam, nameEqualsParam, fuzzyParam)
	response, err := s.sender.SendGetRequest(URL)
	if err != nil {
		return "", err
	}

	json, err := utils.DecodeBytesToMapJson(response)
	if err != nil {
		return "", err
	}

	err = s.ValidateTotalResultsCount(json)
	if err != nil {
		return "", err
	}

	array := json["geonames"].([]interface{})
	obj := array[0].(map[string]interface{})
	name := obj["name"].(string)

	return name, nil
}

func (s *GeoNameService) ValidateCity(city string, country string) (string, models.Coordinates, error) {
	coordinates := models.Coordinates{}

	countryName := countries.ByName(country)
	if countryName == countries.Unknown {
		return "", coordinates, errors.New("country doesn't exist")
	}
	countryCode := countryName.Alpha2()

	featureClassParam := utils.NewHTTPParam("featureClass", "P")
	nameEqualsParam := utils.NewHTTPParam("name_equals", city)
	countryParam := utils.NewHTTPParam("country", countryCode)
	fuzzyParam := utils.NewHTTPParam("fuzzy", s.fuzzy)

	URL := utils.BuildURL(s.preparedURL, featureClassParam, nameEqualsParam, countryParam, fuzzyParam)
	response, err := s.sender.SendGetRequest(URL)
	if err != nil {
		return "", coordinates, errors.New("failing get data from service")
	}

	json, err := utils.DecodeBytesToMapJson(response)
	if err != nil {
		return "", coordinates, err
	}

	err = s.ValidateTotalResultsCount(json)
	if err != nil {
		return "", coordinates, err
	}

	array := json["geonames"].([]interface{})
	obj := array[0].(map[string]interface{})
	name := obj["name"].(string)

	coordinates, err = s.GetCoordinates(json)
	return name, coordinates, err
}

func (s *GeoNameService) ValidateTotalResultsCount(json map[string]interface{}) error {
	totalResultsCount, ok := json["totalResultsCount"].(float64)
	if !ok || totalResultsCount == 0 {
		return errors.New("looks like geo name doesn't exist in our database")
	}
	return nil
}

func (s *GeoNameService) GetCoordinates(json map[string]interface{}) (models.Coordinates, error) {
	coordinates := models.Coordinates{}
	geonamesJson, ok := json["geonames"].([]interface{})

	geonames := geonamesJson[0].(map[string]interface{})
	if !ok {
		return coordinates, errors.New("looks like there is no geo data ")
	}

	lat, ok := geonames["lat"].(string)
	if !ok {
		return coordinates, errors.New("looks like there is no lat for ")
	}

	lon, ok := geonames["lng"].(string)
	if !ok {
		return coordinates, errors.New("looks like there is no lat for ")
	}
	return models.NewCoordinates(lat, lon), nil
}
