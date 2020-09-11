package habitica

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type User struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    struct {
		Auth struct {
			Local struct {
				Username          string `json:"username"`
				LowerCaseUsername string `json:"lowerCaseUsername"`
				Email             string `json:"email"`
			} `json:"local"`
		} `json:"auth"`
		Stats struct {
			Buffs struct {
				Str     int  `json:"str"`
				Int     int  `json:"int"`
				Per     int  `json:"per"`
				Con     int  `json:"con"`
				Stealth int  `json:"stealth"`
				Streaks bool `json:"streaks"`
			} `json:"buffs"`
			Hp     int     `json:"hp"`
			Mp     float64 `json:"mp"`
			Exp    float64 `json:"exp"`
			Gp     float64 `json:"gp"`
			Lvl    int     `json:"lvl"`
			Class  string  `json:"class"`
			Points int     `json:"points"`
			Str    int     `json:"str"`
			Con    int     `json:"con"`
			Int    int     `json:"int"`
			Per    int     `json:"per"`
		} `json:"stats"`
		ID string `json:"id"`
	} `json:"data"`
}

func GetUserInfo() (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.Header.SetMethod("GET")
	SetHeaders(req)
	req.SetRequestURI("https://habitica.com/api/v3/user?userFields=auth.local,stats")

	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Errorf("Failed to make request to Habitica API: %s", err.Error())
		return "", err
	}

	var userData User
	pErr := json.Unmarshal(resp.Body(), &userData)
	if pErr != nil {
		log.Errorf("Failed to parse response: %s", pErr)
		return "", pErr
	}

	if userData.Success {
		return fmt.Sprintf(
			"ID: %s\n"+
				"Username: %s\n"+
				"HP: %d\n"+
				"MP: %.2f",
			userData.Data.ID,
			userData.Data.Auth.Local.Username,
			userData.Data.Stats.Hp,
			userData.Data.Stats.Mp,
		), nil
	} else {
		log.Errorf("Failed to make request to Habitica API: %s", userData.Message)
		return "", fmt.Errorf(userData.Error)
	}
}
