package usecases

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/models"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/services"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/utils"
)

type PredictUsecase struct {
	WeatherService services.WeatherService
}

func (u *PredictUsecase) Handle(update *tgbotapi.Update, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	message := update.Message
	chatID := utils.GetChatId(update)

	_, err := usecaseData.Get(activatedKey)
	if err != nil {
		err := usecaseData.Set(activatedKey, "activated")
		if err != nil {
			return u.getMessageWithInternalError(chatID), Error
		}

		// todo: show current location from DB
		msgCfg := utils.GetMessageWithButtons(chatID, "Do you wanna receive prediction based on your current location?",
			utils.GetInlineButton("Yes", "Please use current location"),
			utils.GetInlineButton("No", "I wanna choose myself"))
		return &msgCfg, Continue
	}

	if update.CallbackQuery != nil {
		return u.HandleCallbacks(update.CallbackQuery, usecaseData)
	} else {
		_, err := usecaseData.Get(isEnterLocationModeSelectedKey)
		if err == nil {

			_, err := usecaseData.Get(countryKey)
			if err != nil {

				err := usecaseData.Set(countryKey, message.Text)
				if err != nil {
					return u.getMessageWithInternalError(chatID), Error
				}

				mes := utils.GetSimpleMessage(chatID, "Enter desired city:\n")
				return &mes, Continue
			} else {
				_, err := usecaseData.Get(cityKey)
				if err != nil {

					err := usecaseData.Set(cityKey, message.Text)
					if err != nil {
						return u.getMessageWithInternalError(chatID), Error
					}

					msgCfg := utils.GetMessageWithButtons(chatID, "Please choose the day you want to receive the weather prediction:",
						utils.GetInlineButton("Next day", "Next day"),
						utils.GetInlineButton("Enter a date", "Enter a date"))
					return &msgCfg, Continue
				} else {
					err := usecaseData.Set(dateKey, message.Text)
					if err != nil {
						return u.getMessageWithInternalError(chatID), Error
					}
					return u.getWeatherMessage(chatID, usecaseData), Finished
				}
			}
		}

		_, err = usecaseData.Get(isEnterDateModeSelectedKey)
		if err == nil {
			err := usecaseData.Set(dateKey, message.Text)
			if err != nil {
				return u.getMessageWithInternalError(chatID), Error
			}
			return u.getWeatherMessage(chatID, usecaseData), Finished
		}
	}

	return u.getMessageWithInternalError(chatID), Finished
}

func (u *PredictUsecase) HandleCallbacks(callbackQuery *tgbotapi.CallbackQuery, usecaseData DataMap) (*tgbotapi.MessageConfig, Status) {
	chatID, callbackData := callbackQuery.Message.Chat.ID, callbackQuery.Data

	switch callbackData {
	case "Please use current location":
		msgCfg := utils.GetMessageWithButtons(chatID, "Please choose the day you want to receive the weather prediction:",
			utils.GetInlineButton("Next day", "Next day"),
			utils.GetInlineButton("Enter a date", "Enter a date"))
		return &msgCfg, Continue
	case "I wanna choose myself":
		_ = usecaseData.Set(isEnterLocationModeSelectedKey, "true")
		msgCfg := utils.GetSimpleMessage(chatID, "Enter desired country:\n")
		return &msgCfg, Continue
	case "Next day":
		return u.getWeatherMessage(chatID, usecaseData), Finished
	case "Enter a date":
		_ = usecaseData.Set(isEnterDateModeSelectedKey, "true")
		msgCfg := utils.GetSimpleMessage(chatID, "Enter desired day:\n")
		return &msgCfg, Continue
	}

	return u.getMessageWithInternalError(chatID), Error
}

func (u *PredictUsecase) mapUsecaseDataToWeatherData(usecaseData DataMap) models.WeatherData {
	country, _ := usecaseData.Get(countryKey)
	city, _ := usecaseData.Get(cityKey)
	date, _ := usecaseData.Get(dateKey)
	return models.WeatherData{Country: country, City: city, Date: utils.MapStringToData(date)}
}

func (u *PredictUsecase) getWeather(usecaseData DataMap) (string, error) {
	return u.WeatherService.GetWeather(u.mapUsecaseDataToWeatherData(usecaseData))
}

func (u *PredictUsecase) getWeatherMessage(chatID int64, usecaseData DataMap) *tgbotapi.MessageConfig {
	weather, err := u.getWeather(usecaseData)
	if err != nil {
		return u.getMessageWithInternalError(chatID)
	}

	msgCfg := tgbotapi.NewMessage(chatID, weather)
	return &msgCfg
}

func (u *PredictUsecase) getMessageWithInternalError(chatID int64) *tgbotapi.MessageConfig {
	invalidMsg := tgbotapi.NewMessage(chatID, "Internal error")
	return &invalidMsg
}
