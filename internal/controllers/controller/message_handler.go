package controller

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/controllers/usecases"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/repositories/temporal_storage"
)

const activeCommandKey = "active_command"

type MessageHandler struct {
	registeredCallbacks map[string]usecases.Usecase
	usecasesData        temporal_storage.TemporalStorage
	bot                 *tgbotapi.BotAPI
}

func NewMessageHandler(bot *tgbotapi.BotAPI) *MessageHandler {
	return &MessageHandler{
		bot: bot,
	}
}

func (h *MessageHandler) RegisterUsecase(usecase usecases.Usecase, command string) error {
	_, ok := h.registeredCallbacks[command]
	if ok {
		return errors.New("callback for this command is already registered")
	}

	h.registeredCallbacks[command] = usecase
	return nil
}

func (h *MessageHandler) ExecuteUsecase(message *tgbotapi.Message, command string) error {
	id := message.Chat.ID
	usecase, exists := h.registeredCallbacks[command]
	if !exists {
		return errors.New("Unrecognised command, also create separate usecase")
	}

	dataAccessor := temporal_storage.CreateTemporalStorageAccessor(h.usecasesData, id)
	mes, status := usecase.Handle(message, &dataAccessor)
	if status == usecases.Finished || status == usecases.Error {
		err := h.usecasesData.Del(id)
		if err != nil {
			// handle
		}
	}

	_, err := h.bot.Send(mes)
	if err != nil {
		// handle
	}

	return nil
}

func extractCommand(text string) (string, error) {
	return "asdf", nil
}

func (h *MessageHandler) AcceptNewMessage(message *tgbotapi.Message) error {
	id := message.Chat.ID

	activeCommand, errActive := h.usecasesData.Get(id, activeCommandKey)
	command, errCommand := extractCommand(message.Text)

	if errActive != nil {
		if errCommand != nil {
			return errors.New("Need to activate separate usecase, ex. /help")
		}
		return h.ExecuteUsecase(message, command)
	}

	if errCommand != nil {
		err := h.usecasesData.Del(id)

		if err != nil {
			// logging
		}
		return h.ExecuteUsecase(message, command)
	}

	return h.ExecuteUsecase(message, activeCommand)
}
