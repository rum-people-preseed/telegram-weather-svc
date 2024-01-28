package upcomming_messages_controller

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/repositories/temporal_storage"
	"github.com/rum-people-preseed/telegram-weather-svc/internal/repositories/user_storage"
)

type UpcommingMessagesController struct {
	temporalStorage temporal_storage.TemporalStorage
	userStorage     user_storage.UserStorage
}

type Command int64

const (
	start   Command = 0
	help    Command = 1
	predict Command = 2
)

func (v *UpcommingMessagesController) HandleMessage(message *tgbotapi.Message) error {
	interactionData, err := v.temporalStorage.GetInteractionData(message.Chat.ID)
	if err == nil {
		return v.handleFollowingInteraction(message, &interactionData)
	}
	return v.handleNewInteraction(message)
}
func getCommand(commandString string) (Command, error) {
	return start, nil
}

func (v *UpcommingMessagesController) handleNewInteraction(message *tgbotapi.Message) error {
	command, err := getCommand(message.Text)
	if err != nil {
		// do something
	}
	print(command)
	return nil
}

func (v *UpcommingMessagesController) handleFollowingInteraction(message *tgbotapi.Message, interactionData *temporal_storage.InteractionData) error {
	return nil
}
