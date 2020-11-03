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
	"todo_web_service/src/telegram_bot/fast_task"
)

const (
	DefaultServiceUrl  = "http://localhost:8080/"
	SuburbanServiceUrl = DefaultServiceUrl + "suburban"
	UserServiceUrl     = DefaultServiceUrl + "user"
)

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

	// В отдельном потоке проверяем fast_task'и.
	go fast_task.CheckFastTasks(&bot)

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

				reply = fmt.Sprintf("Hello %s. This is your 🆔: %s", user.UserName, strconv.Itoa(user.Id))
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
					continue
				}

				err = fast_task.AddFastTask(userId, taskName, chatID, interval)

				if err != nil {
					log.Fatal(err)
				}

				reply = "Задача успешно добавлена!"
			case "/fast_tasks":
				// Получаем все задачи данного пользователя.
				_, reply, err = fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
			case "/delete_fast_task":
				bot.Send(tgbotapi.NewMessage(chatID, "Какая из этих задач уже выполнена? (введите её порядковый номер)"))
				fastTasks, output, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatID, output))

				// Считываем порядковый номер задачи, которую нужно удалить.
				ftUpdate := <-newUpdate
				ftNumber, err := strconv.Atoi(ftUpdate.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Кажется, вы ввели не число. Введите команду ещё раз."))
					continue
				}

				if ftNumber < len(fastTasks)-1 && ftNumber > 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Кажется, такого дела не существует. Введите команду ещё раз."))
					continue
				}

				fastTaskDeleteUrl := DefaultServiceUrl + fmt.Sprintf("%v/fast_task/%v", userId, fastTasks[ftNumber-1].Id)
				_, err = http.NewRequest(http.MethodDelete, fastTaskDeleteUrl, nil)
				if err != nil {
					log.Fatal(err)
				}

				_, output, err = fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Задача %v успешно удалена\n", ftNumber)+output))
			default:
				reply = update.Message.Text
			}

			log.Printf("[%s] - %s", userName, reply)
			msg := tgbotapi.NewMessage(chatID, reply)
			bot.Send(msg)
		}
	}

}
