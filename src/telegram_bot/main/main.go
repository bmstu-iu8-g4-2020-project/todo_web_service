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
	"time"

	"todo_web_service/src/models"
)

const (
	DefaultServiceUrl  = "http://localhost:8080/"
	SuburbanServiceUrl = DefaultServiceUrl + "suburban"
	UserServiceUrl     = DefaultServiceUrl + "user"
)

const (
	emojiAttention = "📢"
)

func CheckFastTasks(bot **tgbotapi.BotAPI) {
	// Содержит время дедлайнов отправки напоминаний о задачах. {id -> time.Time}
	var deadlineTimings map[int]time.Time
	for {
		var allFastTasks []models.FastTask
		resp, err := http.Get(DefaultServiceUrl + "fast_task/")
		if err != nil {
			log.Fatal(err)
		}
		json.NewDecoder(resp.Body).Decode(&allFastTasks)

		// Заполнение дедлайнов.
		for i := range allFastTasks {
			ftId := allFastTasks[i].Id
			// Если время дедлайна нет в мапе, добавляем его.
			if _, inMap := deadlineTimings[ftId]; !inMap {
				deadlineTimings[ftId] = time.Now().Add(allFastTasks[i].NotifyInterval)
			}
		}

		for i := range allFastTasks {
			currFastTask := allFastTasks[i]
			ftId := allFastTasks[i].Id
			// Если дедлайн "просрочен", отправляем напоминание пользователю
			// и обновляем время следующего дедлайна.
			if time.Now().After(deadlineTimings[ftId]) {
				// Чтобы отправить сообщение, нам нужен ChatID...
				(*bot).Send(tgbotapi.NewMessage(currFastTask.ChatId, emojiAttention+currFastTask.TaskName))

				// Увеличиваем дедлайн на величину интервала.
				deadlineTimings[ftId] = deadlineTimings[ftId].Add(allFastTasks[i].NotifyInterval)
			}
		}

		time.Sleep(time.Second * 10)
	}
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
	var userConfig = tgbotapi.NewUpdate(0)

	userConfig.Timeout = 60
	newUpdate, _ := bot.GetUpdatesChan(userConfig)

	go CheckFastTasks(&bot)

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
				user := models.User{
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
			case "/userinfo":
				user := models.User{}

				userInfoUrl := UserServiceUrl + fmt.Sprintf("/%s", strconv.Itoa(userId))

				resp, err := http.Get(userInfoUrl)
				if err != nil {
					log.Fatal(err)
				}

				json.NewDecoder(resp.Body).Decode(&user)

				reply = fmt.Sprintf("Hello %s. This is your id: %s", user.UserName, strconv.Itoa(user.Id))
			case "/suburban":
				resp, err := http.Get(SuburbanServiceUrl)
				if err != nil {
					log.Fatal(err)
				}

				body, _ := ioutil.ReadAll(resp.Body)

				reply = string(body)
			case "/add_fast_task":
				bot.Send(tgbotapi.NewMessage(chatID, "Введите название нового задания."))
				ftUpdate := <-newUpdate
				taskName := ftUpdate.Message.Text

				bot.Send(tgbotapi.NewMessage(chatID, "Введите, с какой периодичностью вам будут приходить сообщения. (Например: 1h10m40s)"))
				ftUpdate = <-newUpdate
				interval, err := time.ParseDuration(ftUpdate.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Кажется, введённое вами сообщение не удовлетворяет формату. Введите команду ещё раз."))
				}

				for ftUpdate = <-newUpdate; err != nil; {
					interval, err = time.ParseDuration(ftUpdate.Message.Text)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(chatID, "Кажется, введённое вами сообщение не удовлетворяет формату. Введите команду ещё раз."))
					}
				}

				fastTask := models.FastTask{
					AssigneeId:     userId,
					TaskName:       taskName,
					ChatId:         chatID,
					NotifyInterval: interval,
					Deadline:       time.Now().Add(interval),
				}

				bytesRepr, err := json.Marshal(fastTask)
				if err != nil {
					log.Fatal(err)
				}

				// DefaultServiceUrl/{id}/fast_task
				fastTaskUrl := DefaultServiceUrl + fmt.Sprintf("/%s", strconv.Itoa(userId)) + "/fast_task"

				_, err = http.Post(fastTaskUrl, "application/json", bytes.NewBuffer(bytesRepr))
				if err != nil {
					log.Fatal(err)
				}

			default:
				reply = update.Message.Text
			}

			log.Printf("[%s] - %s", userName, reply)
			msg := tgbotapi.NewMessage(chatID, reply)
			bot.Send(msg)
		}
	}

}
