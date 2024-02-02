package models

import "time"

type WeatherData struct {
	City        string
	Country     string
	Coordinates Coordinates
	Date        time.Time
}

type Coordinates struct {
	Lat string
	Lon string
}

func NewCoordinates(lat, lon string) Coordinates {
	return Coordinates{
		Lat: lat,
		Lon: lon,
	}
}

func NewWeatherData(city, country string, coordinates Coordinates, date time.Time) WeatherData {
	return WeatherData{
		City:        city,
		Country:     country,
		Coordinates: coordinates,
		Date:        date,
	}
}
