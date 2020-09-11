package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ghotfall/habitica-rebalanced/pkg/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"os"
)

func HandleRequest(_ context.Context, req events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	log.SetFormatter(&log.JSONFormatter{})

	// Bot API
	api, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		log.Errorf("Failed to authorize to the API: %s", err.Error())
		return events.APIGatewayV2HTTPResponse{StatusCode: 403}
	}

	// Request
	u, pErr := bot.ParseUpdate(req.Body)
	if pErr != nil {
		log.Errorf("Failed to parse update: %s", pErr.Error())
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}
	}

	// Processing update
	bot.Entry(api, u)

	return events.APIGatewayV2HTTPResponse{StatusCode: 200}
}

func main() {
	lambda.Start(HandleRequest)
}
