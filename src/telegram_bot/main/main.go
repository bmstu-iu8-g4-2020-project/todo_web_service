package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"todo_web_service/src/models"
	"todo_web_service/src/services"
	"todo_web_service/src/telegram_bot/schedule"

	"todo_web_service/src/telegram_bot/client"
	"todo_web_service/src/telegram_bot/fast_task"
	"todo_web_service/src/telegram_bot/user"

	"github.com/Syfaro/telegram-bot-api"
)

const (
	DefaultServiceUrl  = "http://localhost:8080/"
	SuburbanServiceUrl = DefaultServiceUrl + "suburban"
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
	updates, err := bot.GetUpdatesChan(userConfig)
	if err != nil {
		log.Fatal(err)
	}

	// В отдельном потоке проверяем fast_task'и.
	go fast_task.CheckFastTasks(&bot)

	// читаем обновления из канала
	for update := range updates {
		userName := update.Message.From.UserName
		userId := update.Message.From.ID
		chatID := update.Message.Chat.ID
		var reply string

		switch update.Message.Command() {
		case "start":
			reply, err = user.InitUser(userId, userName,
				update.Message.From.FirstName, update.Message.From.LastName)
			if err != nil {
				log.Fatal(err)
			}
		case "userinfo":
			reply, err = user.GetUserInfo(userId)
			if err != nil {
				log.Fatal(err)
			}
		case "suburban":
			resp, err := http.Get(SuburbanServiceUrl)
			if err != nil {
				log.Fatal(err)
			}

			body, _ := ioutil.ReadAll(resp.Body)

			reply = string(body)
		case "add_fast_task":
			bot.Send(tgbotapi.NewMessage(chatID, "Введите название нового задания."))
			ftUpdate := <-updates
			taskName := ftUpdate.Message.Text

			bot.Send(tgbotapi.NewMessage(chatID,
				"Введите, с какой периодичностью вам будут приходить сообщения. (Например: 1h10m40s)"))
			ftUpdate = <-updates
			interval, err := time.ParseDuration(ftUpdate.Message.Text)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID,
					"Кажется, введённое вами сообщение не удовлетворяет формату. Введите команду ещё раз."))
				continue
			}

			err = fast_task.AddFastTask(userId, taskName, chatID, interval)

			if err != nil {
				log.Fatal(err)
			}

			reply = "Задача успешно добавлена!"
		case "fast_tasks":
			// Получаем все задачи данного пользователя.
			_, reply, err = fast_task.OutputFastTasks(userId)
			if err != nil {
				log.Fatal(err)
			}
		case "delete_fast_task":
			bot.Send(tgbotapi.NewMessage(chatID,
				"Какая из этих задач уже выполнена? (введите её порядковый номер)"))
			fastTasks, output, err := fast_task.OutputFastTasks(userId)
			if err != nil {
				log.Fatal(err)
			}
			bot.Send(tgbotapi.NewMessage(chatID, output))

			// Считываем порядковый номер задачи, которую нужно удалить.
			ftUpdate := <-updates
			ftNumber, err := strconv.Atoi(ftUpdate.Message.Text)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Кажется, вы ввели не число. Введите команду ещё раз."))
				continue
			}

			if ftNumber <= 0 || ftNumber > len(fastTasks) {
				bot.Send(tgbotapi.NewMessage(chatID,
					"Кажется, такого дела не существует. Введите команду ещё раз."))
				continue
			}

			fastTaskDeleteUrl := DefaultServiceUrl +
				fmt.Sprintf("%v/fast_task/%v", userId, fastTasks[ftNumber-1].Id)

			_, err = client.Delete(fastTaskDeleteUrl)

			if err != nil {
				log.Fatal(err)
			}

			_, output, err = fast_task.OutputFastTasks(userId)
			if err != nil {
				log.Fatal(err)
			}

			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Задача %v успешно удалена\n", ftNumber)+output))
		case "fill_schedule":
			bot.Send(tgbotapi.NewMessage(chatID, "Итак, давайте пробежимся по дням недели "+
				"и заполним расписание на каждый из них."))
			// Заполнение всех полей расписания в базе данных расписания.
			assigneeSchedule, err := schedule.InitScheduleTable(userId) // в основном нам нужны отсюда sch_id для заполнения бд.
			if err != nil {
				log.Fatal(err)
			}
			// TODO: Тут должно происходить заполнение полного расписания на неделю.
			var scheduleTasks []models.ScheduleTask // sch_id -> Массив заданий на день.
			for _, weekdaySch := range assigneeSchedule {
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Заполним расписание на %s\n Введите число дел на этот день.",
					services.ParseWeekdayToStr(weekdaySch.WeekDay))))
				// <...>
			}

			err = schedule.FillSchedule(userId, scheduleTasks)
			if err != nil {
				log.Fatal(err)
			}

			bot.Send(tgbotapi.NewMessage(chatID, "Здорово! Ваше расписание успешно заполнено! "))
		default:
			reply = "Введите какую-нибудь команду."
		}

		log.Printf("[%s] - %s", userName, reply)
		msg := tgbotapi.NewMessage(chatID, reply)
		bot.Send(msg)
	}
}
