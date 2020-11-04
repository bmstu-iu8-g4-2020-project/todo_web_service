package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo_web_service/src/models"
)

const (
	DefaultServiceUrl = "http://localhost:8080/"
	UserServiceUrl    = DefaultServiceUrl + "user"
)

func InitUser(userId int, userName string, firstName string, secondName string) (string, error) {
	reply := fmt.Sprintf("Здравствуйте, %s %s!\n", firstName, secondName)

	user := models.User{
		Id:         userId,
		UserName:   userName,
		FirstName:  firstName,
		SecondName: secondName,
	}

	bytesRepr, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	_, err = http.Post(UserServiceUrl, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		return "", err
	}

	reply += "Добро пожаловать!"

	return reply, nil
}

func GetUserInfo(userId int) (string, error) {
	user := models.User{}

	userInfoUrl := UserServiceUrl + fmt.Sprintf("/%s", strconv.Itoa(userId))

	resp, err := http.Get(userInfoUrl)
	if err != nil {
		return "", err
	}

	json.NewDecoder(resp.Body).Decode(&user)

	reply := fmt.Sprintf("Здравствуйте, %s %s. \nВаш 🆔: %s",
		user.FirstName, user.SecondName, strconv.Itoa(user.Id))

	return reply, nil
}
