package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todo_web_service/src/models"
	"todo_web_service/src/telegram_bot/utils"
)

const (
	DefaultServiceUrl = "http://localhost:8080/"
)

type State struct {
	Code    int    `json:"code"`
	Request string `json:"request"`
}

func InitUser(userId int, userName string) error {
	user := models.User{
		Id:           userId,
		UserName:     userName,
		StateCode:    START,
		StateRequest: "{}",
	}

	bytesRepr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	url := DefaultServiceUrl + "user/"
	_, err = http.Post(url, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		return err
	}

	return nil
}

func GetUser(userId int) (models.User, error) {
	user := models.User{}

	url := DefaultServiceUrl + fmt.Sprintf("user/%v", userId)

	resp, err := http.Get(url)
	if err != nil {
		return models.User{}, err
	}

	json.NewDecoder(resp.Body).Decode(&user)

	return user, nil
}

func GetStates(userStates *map[int]State) error {
	url := DefaultServiceUrl + "user/"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	var users []models.User
	json.NewDecoder(resp.Body).Decode(&users)

	for _, user := range users {
		(*userStates)[user.Id] = State{user.StateCode, user.StateRequest}
	}

	return nil
}

func UpdateUser(userId int, username string, stateCode int, stateRequest string) error {
	user := models.User{
		Id:           userId,
		UserName:     username,
		StateCode:    stateCode,
		StateRequest: stateRequest,
	}
	url := DefaultServiceUrl + "user/{id}"

	bytesRepr, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = utils.Put(url, bytes.NewBuffer(bytesRepr))
	if err != nil {
		return err
	}

	return err
}

func SetState(userId int, userName string, userStates *map[int]State, state State) error {
	err := UpdateUser(userId, userName, state.Code, state.Request)
	if err != nil {
		return err
	}
	(*userStates)[userId] = State{Code: state.Code, Request: state.Request}

	return nil
}

func ResetState(userId int, userName string, userStates *map[int]State) error {
	err := UpdateUser(userId, userName, START, "{}")
	if err != nil {
		return err
	}
	(*userStates)[userId] = State{Code: START, Request: "{}"}

	return nil
}
