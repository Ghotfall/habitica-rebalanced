package bot

import (
	"encoding/json"
	"github.com/ghotfall/habitica-rebalanced/pkg/habitica"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func ParseUpdate(data string) (*tgbotapi.Update, error) {
	var u tgbotapi.Update
	err := json.Unmarshal([]byte(data), &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func Entry(api *tgbotapi.BotAPI, u *tgbotapi.Update) {
	if u.CallbackQuery != nil {
		resp, err := habitica.GetUserTask(u.CallbackQuery.Data)
		if err == nil && u.CallbackQuery.Message != nil {
			_, mErr := api.Send(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, resp))
			if mErr != nil {
				log.Errorf("Failed to send answer to user: %s", mErr.Error())
			}

			_, cErr := api.AnswerCallbackQuery(tgbotapi.CallbackConfig{
				CallbackQueryID: u.CallbackQuery.ID,
				Text:            "Processing...",
			})
			if cErr != nil {
				log.Errorf("Failed to send callback answer to user: %s", cErr.Error())
			}
		}

	} else if u.Message != nil {
		if u.Message.IsCommand() {
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")

			switch u.Message.Command() {
			case "get_user":
				resp, err := habitica.GetUserInfo()
				if err == nil {
					msg.Text = resp
				}

			case "get_tasks":
				resp, err := habitica.GetUserTasks()
				if err == nil {
					msg.Text = "Your tasks:"
					msg.ReplyMarkup = resp
				}
			}

			if msg.Text != "" {
				_, err := api.Send(msg)
				if err != nil {
					log.Errorf("Failed to send answer to user: %s", err.Error())
				}
			}
		}
	}
}
