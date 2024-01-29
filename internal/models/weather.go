package models

import "time"

type WeatherData struct {
	Country string
	City    string
	Date    time.Time
}
