package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"os"
)

func HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.SetFormatter(&log.JSONFormatter{})
	var u tgbotapi.Update
	if err := json.Unmarshal([]byte(req.Body), &u); err == nil {
		if u.Message == nil {
			return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
		}

		if !u.Message.IsCommand() {
			return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
		}

		bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
		if err != nil {
			return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
		}

		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
		switch u.Message.Command() {
		case "start":
			msg.Text = "Hello! I'm really working now 🤖"
		case "help":
			msg.Text = "Sorry, no functionality has been implemented yet 🙃"
		case "status":
			msg.Text = "I'm ok 🤖"
		case "dice":
			go bot.Send(tgbotapi.NewDice(u.Message.Chat.ID))
			return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
		case "open":
			msg.Text = u.Message.Command()
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("A"),
					tgbotapi.NewKeyboardButton("B"),
					tgbotapi.NewKeyboardButton("C"),
				))
		case "close":
			msg.Text = u.Message.Command()
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "whoami":
			switch u.Message.CommandArguments() {
			case "aws":
				lc, _ := lambdacontext.FromContext(ctx)
				msg.Text = lc.InvokedFunctionArn
			default:
				msg.Text = u.Message.From.UserName
			}
		default:
			msg.Text = fmt.Sprintf("I don't know this command: `%s`", u.Message.Command())
		}

		go bot.Send(msg)

	} else {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, fmt.Errorf("failed to unmarshall object")
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
