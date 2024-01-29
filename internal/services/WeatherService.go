package services

import (
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
)

type (
	WeatherService interface {
		GetWeather(weatherData models.WeatherData) (string, error)
	}

	WeatherProvider struct {
	}
)

func (w *WeatherProvider) GetWeather(weatherData models.WeatherData) (string, error) {
	// todo: send request to service. show received data.
	return "The next day will be hot just like you!", nil
}
