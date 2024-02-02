package date_parser

import (
	"errors"
	"github.com/araddon/dateparse"
	"time"
)

type DateParser interface {
	ParseDateString(dateString string) (time.Time, error)
}

type MultiformatDateParser struct {
}

func (p *MultiformatDateParser) ParseDateString(dateString string) (time.Time, error) {
	parsedTime, err := dateparse.ParseAny(dateString)
	if err != nil {
		return p.parseNearlyDate(dateString)
	}

	return parsedTime, nil
}

func (p *MultiformatDateParser) parseNearlyDate(dateString string) (time.Time, error) {
	return time.Time{}, errors.New("failed to parse date")
}
