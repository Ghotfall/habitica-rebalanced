package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

func main() {
	lambda.Start(handle)
}

func handle(update tgbotapi.Update) error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		return err
	}

	bot.Debug = true

	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//if _, err := bot.Send(msg); err != nil {
	//	return err
	//}

	fmt.Printf("%+v", update)

	return nil
}
