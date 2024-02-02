package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/biter777/countries"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"io"
	"net/http"
	"os"
	"strings"
)

type GeoService interface {
	ValidateCountry(country string) error
	ValidateCity(city string, country string) error
}

type GeoNameService struct {
	baseURL string
	log     controller.Logger
}

func NewGeoNameService(logger controller.Logger) GeoNameService {
	baseURL := "http://api.geonames.org/searchJSON?maxRows=1&username=%v"
	geoNameUsername := os.Getenv("GEO_NAME_SERVICE_USERNAME")
	return GeoNameService{baseURL: fmt.Sprintf(baseURL, geoNameUsername), log: logger}
}

func (s *GeoNameService) ValidateCountry(country string) error {
	jsonResult, err := s.SendGetRequestWithParams("featureClass=A", "name_equals="+country)
	if err != nil {
		return err
	}

	_, err = s.GetTotalResultsCount(jsonResult)
	if err != nil {
		return err
	}

	// here we can already return coordinates for city/country/etc
	return nil
}

func (s *GeoNameService) ValidateCity(city string, country string) error {

	countryName := countries.ByName(country)
	if countryName == countries.Unknown {
		return errors.New("country doesn't exist")
	}

	countryCode := countryName.Alpha2()
	jsonResult, err := s.SendGetRequestWithParams("featureClass=P", "name_equals="+city, "country="+countryCode)
	if err != nil {
		return err
	}

	_, err = s.GetTotalResultsCount(jsonResult)
	if err != nil {
		return err
	}

	// here we can already return coordinates for city/country/etc
	return nil
}

func (s *GeoNameService) GetUrlWithQueryParam(params ...string) string {
	queryValue := strings.Join(params, "&")
	return strings.Join([]string{s.baseURL, queryValue}, "&")
}

func (s *GeoNameService) SendGetRequestWithParams(params ...string) (map[string]interface{}, error) {
	queryURL := s.GetUrlWithQueryParam(params...)
	resp, err := http.Get(queryURL)
	s.log.Infof("URl to Geodata service is sent. URL - " + queryURL)

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

func (s *GeoNameService) GetTotalResultsCount(json map[string]interface{}) (float64, error) {
	totalResultsCount, ok := json["totalResultsCount"].(float64)
	if !ok || totalResultsCount == 0 {
		fmt.Println("Error parsing totalResultsCount")
		return 0, errors.New("looks like geo name don't exist in our database")
	}
	return totalResultsCount, nil
}
