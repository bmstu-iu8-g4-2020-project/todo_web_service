package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	DefaultServiceUrl  = "http://localhost:8080/"
	SuburbanServiceUrl = DefaultServiceUrl + "suburban"
	UserServiceUrl     = DefaultServiceUrl + "user"
)

type Message struct {
	UserName string `json:"username"`
	ChatID   int64  `json:"chat_id"`
	Text     string `json:"text"`
}

type User struct {
	Id       int    `json:"user_id"`
	UserName string `json:"username"`
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	_, err = bot.RemoveWebhook()
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var userConfig tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)

	userConfig.Timeout = 60
	newUpdate, _ := bot.GetUpdatesChan(userConfig)

	// читаем обновления из канала
	for {
		select {
		case update := <-newUpdate: //  получить из канала
			userName := update.Message.From.UserName
			userId := update.Message.From.ID
			chatID := update.Message.Chat.ID
			var reply string

			switch update.Message.Text {
			case "/start":
				reply = fmt.Sprintf("Hello %s!\n Welcome =)", userName)
				user := User{
					Id:       userId,
					UserName: userName,
				}

				bytesRepr, err := json.Marshal(user)
				if err != nil {
					log.Fatal(err)
				}

				_, err = http.Post(UserServiceUrl, "application/json", bytes.NewBuffer(bytesRepr))
				if err != nil {
					log.Fatal(err)
				}

				reply += fmt.Sprintf("\nВы авторизованы!")
			case "/services":
				msg := Message{
					UserName: userName,
					ChatID:   chatID,
				}

				bytesRepr, err := json.Marshal(msg)
				if err != nil {
					log.Fatal(err)
				}

				resp, err := http.Post(DefaultServiceUrl, "application/json", bytes.NewBuffer(bytesRepr))
				if err != nil {
					log.Fatal(err)
				}

				json.NewDecoder(resp.Body).Decode(&msg)

				reply = msg.Text
			case "/userinfo":
				user := User{}

				userInfoUrl := UserServiceUrl + fmt.Sprintf("/%s", strconv.Itoa(userId))

				resp, err := http.Get(userInfoUrl)
				if err != nil {
					log.Fatal(err)
				}

				json.NewDecoder(resp.Body).Decode(&user)

				reply = fmt.Sprintf("Hello %s. This is your id: %s", user.UserName, user.Id)
			case "/suburban":
				resp, err := http.Get(SuburbanServiceUrl)
				if err != nil {
					log.Fatal(err)
				}

				body, _ := ioutil.ReadAll(resp.Body)

				reply = string(body)
			default:
				reply = update.Message.Text
			}

			log.Printf("[%s] - %s", userName, reply)
			msg := tgbotapi.NewMessage(chatID, reply)
			bot.Send(msg)
		}
	}

}
