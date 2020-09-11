package bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	if u.Message != nil {
		api.Send(tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text))
	}
}
