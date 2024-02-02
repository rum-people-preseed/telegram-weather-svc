package services

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/biter777/countries"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services/utils"

	"io"
	"net/http"
)

type GeoService interface {
	ValidateCountry(country string) error
	ValidateCity(city string, country string) (string, string, error)
}

type GeoNameService struct {
	preparedURL string
	log         models.Logger
}

func NewGeoNameService(logger models.Logger) GeoNameService {
	baseURL := "http://api.geonames.org/searchJSON"
	apiKey := os.Getenv("GEO_NAME_SERVICE_USERNAME")
	apiKeyParam := utils.NewHTTPParam("username", apiKey)
	maxRowsParam := utils.NewHTTPParam("maxRows", "1")
	preparedURL := utils.BuildURL(baseURL, apiKeyParam, maxRowsParam)
	return GeoNameService{preparedURL: preparedURL, log: logger}
}

func (s *GeoNameService) ValidateCountry(country string) error {
	featureClassParam := utils.NewHTTPParam("featureClass", "A")
	nameEqualsParam := utils.NewHTTPParam("name_equals", country)

	jsonResult, err := s.SendGetRequestWithParams(featureClassParam, nameEqualsParam)
	if err != nil {
		return err
	}

	err = s.ValidateTotalResultsCount(jsonResult)
	if err != nil {
		return err
	}

	// here we can already return coordinates for city/country/etc
	return nil
}

func (s *GeoNameService) ValidateCity(city string, country string) (string, string, error) {
	lat, lon := "", ""

	countryName := countries.ByName(country)
	if countryName == countries.Unknown {
		return lat, lon, errors.New("country doesn't exist")
	}
	countryCode := countryName.Alpha2()

	featureClassParam := utils.NewHTTPParam("featureClass", "P")
	nameEqualsParam := utils.NewHTTPParam("name_equals", city)
	countryParam := utils.NewHTTPParam("country", countryCode)

	jsonResult, err := s.SendGetRequestWithParams(featureClassParam, nameEqualsParam, countryParam)
	if err != nil {
		return lat, lon, err
	}

	err = s.ValidateTotalResultsCount(jsonResult)
	if err != nil {
		return lat, lon, err
	}

	// here we can already return coordinates for city/country/etc
	lat, lon, err = s.GetCoordinates(jsonResult)
	return lat, lon, err
}

func (s *GeoNameService) SendGetRequestWithParams(params ...*utils.HTTPParam) (map[string]interface{}, error) {
	url := utils.BuildURL(s.preparedURL, params...)
	resp, err := http.Get(url)
	s.log.Infof("URl to Geodata service is sent. URL - " + url)

	if err != nil {
		return nil, errors.New("failing get info from service")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failing read response from service")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("error decoding response from service")
	}

	return result, nil
}

func (s *GeoNameService) ValidateTotalResultsCount(json map[string]interface{}) error {
	totalResultsCount, ok := json["totalResultsCount"].(float64)
	if !ok || totalResultsCount == 0 {
		return errors.New("looks like geo name doesn't exist in our database")
	}
	return nil
}

func (s *GeoNameService) GetCoordinates(json map[string]interface{}) (string, string, error) {
	geonamesJson, ok := json["geonames"].([]interface{})

	geonames := geonamesJson[0].(map[string]interface{})
	if !ok {
		return "", "", errors.New("looks like there is no geo data ")
	}

	lat, ok := geonames["lat"].(string)
	if !ok {
		return "", "", errors.New("looks like there is no lat for ")
	}

	lon, ok := geonames["lng"].(string)
	if !ok {
		return "", "", errors.New("looks like there is no lat for ")
	}
	return lat, lon, nil
}
