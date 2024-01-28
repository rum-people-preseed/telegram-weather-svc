package temporal_storage

type InteractionType int64

const (
	GettingLocationData InteractionType = 0
	GettingDataTime     InteractionType = 1
)

type InteractionData struct {
	chatId          int64
	interactionType InteractionType
	city            string
	country         string
}

type TemporalStorage interface {
	AddNewInteraction(chatId int64, interactionType InteractionType) error
	AddCountry(chatId int64, country string) error
	AddCity(chatId int64, city string) error
	GetInteractionData(chatId int64) (InteractionData, error)
}
