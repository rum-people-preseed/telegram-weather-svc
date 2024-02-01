package usecases

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/controller"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/sequences"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
)

type UpdateLocationUsecaseFactory struct {
	GeoService services.GeoService
}

func (f *UpdateLocationUsecaseFactory) Create() controller.Usecase {
	return &UpdateLocationUsecase{
		locationSequence: sequences.CreateGetLocationSequence(f.GeoService),
	}
}

func (f *UpdateLocationUsecaseFactory) Command() string {
	return "/update_location"
}

type UpdateLocationUsecase struct {
	locationSequence sequences.GetLocationSequence
	state            string
}

func (u *UpdateLocationUsecase) Handle(update *tgbotapi.Update) (*tgbotapi.MessageConfig, controller.Status) {
	return u.HandleInitialState(update)
}

func (u *UpdateLocationUsecase) HandleInitialState(update *tgbotapi.Update) (*tgbotapi.MessageConfig, controller.Status) {
	mes, err := u.locationSequence.Handle(update)

	if err == controller.Continue {
		return mes, err
	}
	if err == controller.Error {
		// todo handle error state
	}

	message := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
		"You location has been updated: %v, %v",
		u.locationSequence.GetCountryName(), u.locationSequence.GetCityName()))

	return &message, controller.Finished
}
