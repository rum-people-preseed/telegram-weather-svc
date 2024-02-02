package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services/utils"
)

type (
	WeatherService interface {
		GetWeatherPredictionMessage(weatherData models.WeatherData, chatID int64) (tgbotapi.Chattable, error)
	}

	WeatherPredictorService struct {
		URL string
		log models.Logger
	}
)

func NewWeatherPredictorService(URL string, log models.Logger) WeatherPredictorService {
	return WeatherPredictorService{URL: URL, log: log}
}

func (w *WeatherPredictorService) GetWeatherPredictionMessage(weatherData models.WeatherData, chatID int64) (tgbotapi.Chattable, error) {
	avrgTemp, chartBase64, err := w.GetTemperatureAndChart(weatherData)
	if err != nil {
		return nil, err
	}

	chartBytes, err := base64.StdEncoding.DecodeString(chartBase64)
	if err != nil {
		return nil, errors.New("error during decoding chart")
	}

	imageBuffer := bytes.NewBuffer(chartBytes)
	image := tgbotapi.FileBytes{Name: "chart.png", Bytes: imageBuffer.Bytes()}

	photoConfig := tgbotapi.NewPhotoUpload(chatID, image)
	photoConfig.Caption = fmt.Sprintf("Average temperature is %v°", avrgTemp)
	return photoConfig, nil
}

func (w *WeatherPredictorService) GetParamsFromWeatherData(weatherData models.WeatherData) []*utils.HTTPParam {
	params := []*utils.HTTPParam{
		utils.NewHTTPParam("city", weatherData.City),
		utils.NewHTTPParam("country", weatherData.Country),
		utils.NewHTTPParam("date", weatherData.Date),
		utils.NewHTTPParam("latitude", weatherData.Lat),
		utils.NewHTTPParam("longtitude", weatherData.Lon),
	}
	return params
}

func (w *WeatherPredictorService) GetTemperatureAndChart(weatherData models.WeatherData) (float64, string, error) {
	avrgTemp, chartBase64 := 0.0, ""

	svcResponse, err := w.SendGetRequestWithParams(weatherData)
	if err != nil {
		w.log.Error(err)
		return avrgTemp, chartBase64, errors.New("failing get info from service")
	}

	avrgTemp, ok := svcResponse["average_temperature"].(float64)
	if !ok {
		w.log.Error("there is no avrg temperature from response")
		return avrgTemp, chartBase64, errors.New("there is no avrg temperature from response")
	}

	chartBase64, ok = svcResponse["chart"].(string)
	if !ok {
		w.log.Error("there is no chart from response")
		return avrgTemp, chartBase64, errors.New("there is no chart from response")
	}

	return avrgTemp, chartBase64, nil
}

func (w *WeatherPredictorService) SendGetRequestWithParams(weatherData models.WeatherData) (map[string]interface{}, error) {
	url := utils.BuildURL(w.URL+"/temperature/", w.GetParamsFromWeatherData(weatherData)...)
	w.log.Infof("URl to weather service service is sent. URL - " + url)
	resp, err := http.Get(url)

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
