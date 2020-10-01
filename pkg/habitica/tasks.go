package habitica

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strings"
)

type Tasks struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    []struct {
		ID    string `json:"_id"`
		Notes string `json:"notes"`
		Text  string `json:"text"`
	} `json:"data"`
}

type Task struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    struct {
		ID    string `json:"_id"`
		Notes string `json:"notes"`
		Text  string `json:"text"`
	} `json:"data"`
}

func (t *Tasks) String() string {
	var str strings.Builder
	for i, task := range t.Data {
		_, err := fmt.Fprintf(&str, "%d) %s\n", i+1, task.Text)
		if err != nil {
			continue
		}
	}

	return str.String()
}

func (t *Tasks) InlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(t.Data))

	for _, task := range t.Data {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(task.Text, task.ID),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func GetUserTasks() (tgbotapi.InlineKeyboardMarkup, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.Header.SetMethod("GET")
	SetHeaders(req)
	req.SetRequestURI("https://habitica.com/api/v3/tasks/user?type=todos")

	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Errorf("Failed to make request to Habitica API: %s", err.Error())
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	var userTasks Tasks
	pErr := json.Unmarshal(resp.Body(), &userTasks)
	if pErr != nil {
		log.Errorf("Failed to parse response: %s", pErr)
		return tgbotapi.InlineKeyboardMarkup{}, pErr
	}

	if userTasks.Success {
		return userTasks.InlineKeyboard(), nil
	} else {
		log.Errorf("Failed to make request to Habitica API: %s", userTasks.Message)
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf(userTasks.Error)
	}
}

func GetUserTask(id string) (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.Header.SetMethod("GET")
	SetHeaders(req)
	req.SetRequestURI("https://habitica.com/api/v3/tasks/" + id)

	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Errorf("Failed to make request to Habitica API: %s", err.Error())
		return "", err
	}

	var userTask Task
	pErr := json.Unmarshal(resp.Body(), &userTask)
	if pErr != nil {
		log.Errorf("Failed to parse response: %s", pErr)
		return "", pErr
	}

	if userTask.Success {
		return userTask.Data.Text, nil
	} else {
		log.Errorf("Failed to make request to Habitica API: %s", userTask.Message)
		return "", fmt.Errorf(userTask.Error)
	}
}
