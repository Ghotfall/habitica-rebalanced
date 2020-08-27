package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	lambda.Start(handle)
}

func handle(_ context.Context, update tgbotapi.Update) (string, error) {
	if update.Message != nil {
		return fmt.Sprintf("Message: %s", update.Message.Text), nil
	} else {
		return "Not a text message", nil
	}
}
