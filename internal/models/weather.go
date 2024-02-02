package models

import "time"

type WeatherData struct {
	City    string
	Country string
	Lat     string
	Lon     string
	Date    string
}

func NewWeatherData(city, country, lat, lon, date string) WeatherData {
	return WeatherData{
		City:    city,
		Country: country,
		Lat:     lat,
		Lon:     lon,
		Date:    GetFormattedDate(date),
	}
}

func GetFormattedDate(date string) string {
	return time.Now().Format("2016-12-30 09:00:00")
}
