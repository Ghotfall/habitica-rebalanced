package bot

import (
	"encoding/json"
	"github.com/ghotfall/habitica-rebalanced/pkg/habitica"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strings"
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
		// Split data string to get exact command
		spl := strings.Split(u.CallbackQuery.Data, " ")

		// Find this command
		switch spl[0] {
		case "get":
			GetTask(api, spl[1], u.CallbackQuery)
		}

	} else if u.Message != nil {
		if u.Message.IsCommand() {
			// Create new message
			msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")

			// Fort message body depending on used command
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

			// Send message if body is not empty
			if msg.Text != "" {
				_, err := api.Send(msg)
				if err != nil {
					log.Errorf("Failed to send answer to user: %s", err.Error())
				}
			}
		}
	}
}
