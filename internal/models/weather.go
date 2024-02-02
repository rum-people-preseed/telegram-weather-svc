package models

import "time"

type WeatherData struct {
	City    string
	Country string
	Lat     string
	Lon     string
	Date    time.Time
}

func NewWeatherData(city, country, lat, lon string, date time.Time) WeatherData {
	return WeatherData{
		City:    city,
		Country: country,
		Lat:     lat,
		Lon:     lon,
		Date:    date,
	}
}
