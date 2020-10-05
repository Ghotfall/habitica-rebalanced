package bot

import (
	"github.com/ghotfall/habitica-rebalanced/pkg/habitica"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// Later: change names for habitica and bot functions
func GetTask(api *tgbotapi.BotAPI, data string, query *tgbotapi.CallbackQuery) {
	resp, err := habitica.GetUserTask(data)
	if err == nil && query.Message != nil {
		_, mErr := api.Send(tgbotapi.NewMessage(query.Message.Chat.ID, resp))
		if mErr != nil {
			log.Errorf("Failed to send answer to user: %s", mErr.Error())
		}

		_, cErr := api.AnswerCallbackQuery(tgbotapi.CallbackConfig{
			CallbackQueryID: query.ID,
			Text:            "Processing...",
		})
		if cErr != nil {
			log.Errorf("Failed to send callback answer to user: %s", cErr.Error())
		}
	}
}
