package habitica

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strings"
)

type Tasks struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    []struct {
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

func GetUserTasks() (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.Header.SetMethod("GET")
	SetHeaders(req)
	req.SetRequestURI("https://habitica.com/api/v3/tasks/user?type=todos")

	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Errorf("Failed to make request to Habitica API: %s", err.Error())
		return "", err
	}

	var userTasks Tasks
	pErr := json.Unmarshal(resp.Body(), &userTasks)
	if pErr != nil {
		log.Errorf("Failed to parse response: %s", pErr)
		return "", pErr
	}

	if userTasks.Success {
		return userTasks.String(), nil
	} else {
		log.Errorf("Failed to make request to Habitica API: %s", userTasks.Message)
		return "", fmt.Errorf(userTasks.Error)
	}
}
