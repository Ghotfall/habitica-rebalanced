package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/ghotfall/habitica-rebalanced/pkg/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.SetFormatter(&log.JSONFormatter{})

	// Request
	u, pErr := bot.ParseUpdate(req.Body)
	if pErr != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, pErr
	}

	// Bot API
	api, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	// Processing update
	if u.CallbackQuery != nil {
		c, err := processCallback(u.CallbackQuery)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
		}

		sendMsg(api, c)

	} else if u.Message.IsCommand() {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
		switch u.Message.Command() {
		case "start":
			msg.Text = "Hello! I'm really working now 🤖"
		case "help":
			msg.Text = "Sorry, no functionality has been implemented yet 🙃"
		case "status":
			msg.Text = "I'm ok 🤖"
		case "dice":
			go api.Send(tgbotapi.NewDice(u.Message.Chat.ID))
			return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
		case "score":
			msg.Text = "Score: 0"
			msg.ReplyMarkup = getKeyboardMarkup()
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

		sendMsg(api, msg)
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 200}, nil
}

func sendMsg(api *tgbotapi.BotAPI, c tgbotapi.Chattable) {
	_, err := api.Send(c)
	if err != nil {
		log.Errorf("Failed to send message: %s", err.Error())
	}
}

func processCallback(query *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	if query.Message != nil {
		scoreStr := strings.Fields(query.Message.Text)
		scoreInt, err := strconv.Atoi(scoreStr[1])
		if err != nil {
			log.Error("Internal error: failed to parse score")
			return nil, err
		}

		switch query.Data {
		case "score_p1":
			scoreInt++
		case "score_m1":
			scoreInt--
		case "score_r":
			scoreInt = 0
		}
		msg := fmt.Sprintf("Score: %d", scoreInt)
		return tgbotapi.NewEditMessageTextAndMarkup(
			query.Message.Chat.ID,
			query.Message.MessageID,
			msg,
			getKeyboardMarkup()), nil
	} else {
		return nil, fmt.Errorf("callback message is empty")
	}
}

func getKeyboardMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+1", "score_p1"),
			tgbotapi.NewInlineKeyboardButtonData("-1", "score_m1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Reset", "score_r"),
		),
	)
}

type ReplyPayload struct {
	Method string `json:"method"`
	Params url.Values
}

func main() {
	lambda.Start(HandleRequest)
}
