package utils

import (
	"time"
)

func GetNextDay() string {
	nextDay := time.Now().AddDate(0, 0, 1)
	return nextDay.Format("02.01.2006")
}

func MapStringToData(data string) time.Time {
	return time.Now()
}
