package services

import (
	"errors"
	"io"
	"net/http"

	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
)

type Sender interface {
	SendGetRequest(URL string) ([]byte, error)
}

type HHTPSender struct {
	Log models.Logger
}

func (s *HHTPSender) SendGetRequest(URL string) ([]byte, error) {
	s.Log.Infof("Sending %v", URL)
	resp, err := http.Get(URL)

	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("failing get info from service")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Log.Error(err)
		return nil, errors.New("failing read response from service")
	}

	return body, err
}
