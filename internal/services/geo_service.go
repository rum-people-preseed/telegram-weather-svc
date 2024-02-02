package services

import (
	"errors"
	"os"

	"github.com/biter777/countries"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services/utils"
)

type GeoService interface {
	ValidateCountry(country string) error
	ValidateCity(city string, country string) (models.Coordinates, error)
}

type GeoNameService struct {
	preparedURL string
	sender      Sender
	log         models.Logger
}

func NewGeoNameService(sender Sender, logger models.Logger) GeoNameService {
	baseURL := "http://api.geonames.org/searchJSON"
	apiKey := os.Getenv("GEO_NAME_SERVICE_USERNAME")
	apiKeyParam := utils.NewHTTPParam("username", apiKey)
	maxRowsParam := utils.NewHTTPParam("maxRows", "1")
	preparedURL := utils.BuildURL(baseURL, apiKeyParam, maxRowsParam)
	return GeoNameService{preparedURL: preparedURL, sender: sender, log: logger}
}

func (s *GeoNameService) ValidateCountry(country string) error {
	featureClassParam := utils.NewHTTPParam("featureClass", "A")
	nameEqualsParam := utils.NewHTTPParam("name_equals", country)

	URL := utils.BuildURL(s.preparedURL, featureClassParam, nameEqualsParam)
	response, err := s.sender.SendGetRequest(URL)
	if err != nil {
		return errors.Join(err, errors.New("failing get data from service"))
	}

	json, err := utils.DecodeBytesToMapJson(response)
	if err != nil {
		return errors.Join(err, errors.New("failing get data from service"))
	}

	err = s.ValidateTotalResultsCount(json)
	if err != nil {
		return err
	}

	return nil
}

func (s *GeoNameService) ValidateCity(city string, country string) (models.Coordinates, error) {
	coordinates := models.Coordinates{}

	countryName := countries.ByName(country)
	if countryName == countries.Unknown {
		return coordinates, errors.New("country doesn't exist")
	}
	countryCode := countryName.Alpha2()

	featureClassParam := utils.NewHTTPParam("featureClass", "P")
	nameEqualsParam := utils.NewHTTPParam("name_equals", city)
	countryParam := utils.NewHTTPParam("country", countryCode)

	URL := utils.BuildURL(s.preparedURL, featureClassParam, nameEqualsParam, countryParam)
	response, err := s.sender.SendGetRequest(URL)
	if err != nil {
		return coordinates, errors.New("failing get data from service")
	}

	json, err := utils.DecodeBytesToMapJson(response)
	if err != nil {
		return coordinates, err
	}

	err = s.ValidateTotalResultsCount(json)
	if err != nil {
		return coordinates, err
	}

	coordinates, err = s.GetCoordinates(json)
	return coordinates, err
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
