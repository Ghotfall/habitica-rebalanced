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
