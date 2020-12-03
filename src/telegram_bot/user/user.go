package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/models"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/utils"
)

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

	url := utils.DefaultServiceUrl + "user/"
	_, err = http.Post(url, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		return err
	}

	return nil
}

func GetUser(userId int) (models.User, error) {
	user := models.User{}

	url := utils.DefaultServiceUrl + fmt.Sprintf("user/%v", userId)

	resp, err := http.Get(url)
	if err != nil {
		return models.User{}, err
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func UpdateUser(userId int, username string, stateCode int, stateRequest string) error {
	user := models.User{
		Id:           userId,
		UserName:     username,
		StateCode:    stateCode,
		StateRequest: stateRequest,
	}

	bytesRepr, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = utils.Put(utils.DefaultServiceUrl+"user/{id}", bytes.NewBuffer(bytesRepr))
	if err != nil {
		return err
	}

	return err
}
