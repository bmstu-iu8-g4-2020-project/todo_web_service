package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"todo_web_service/src/models"

	"todo_web_service/src/telegram_bot/fast_task"
	"todo_web_service/src/telegram_bot/user"
	"todo_web_service/src/telegram_bot/utils"

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

	// Храним состояния пользователей в бд.
	userStates := make(map[int]user.State)
	err = user.GetStates(&userStates)
	if err != nil {
		log.Fatal(err)
	}
	// читаем обновления из канала
	for update := range updates {
		chatID := update.Message.Chat.ID
		msg := update.Message
		userId := update.Message.From.ID
		userName := update.Message.From.UserName

		switch update.Message.Command() {
		case "start":
			if userStates[userId].Code == user.START {
				err = user.InitUser(userId, userName)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Здравствуйте, %s.\nДобро пожаловать!", userName)))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Вы не закончили ввод данных."))
			}
			continue
		case "userinfo":
			if userStates[userId].Code == user.START {
				user, err := user.GetUser(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Здравствуйте, %s. \nВаш 🆔: %s",
					user.UserName, strconv.Itoa(user.Id))))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Вы не закончили ввод данных."))
			}
			continue
		case "suburban":
			if userStates[userId].Code == user.START {
				resp, err := http.Get(SuburbanServiceUrl)
				if err != nil {
					log.Fatal(err)
				}

				body, _ := ioutil.ReadAll(resp.Body)

				bot.Send(tgbotapi.NewMessage(chatID, string(body)))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Вы не закончили ввод данных."))
			}
			continue
		case "add_fast_task":
			if userStates[userId].Code == user.START {
				state := user.State{Code: user.FAST_TASK_ENTER_TITLE, Request: "{}"}
				user.SetState(userId, userName, &userStates, state)
				bot.Send(tgbotapi.NewMessage(chatID, "Введите название нового задания."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Вы не закончили ввод данных."))
			}
			continue
		case "fast_tasks":
			if userStates[userId].Code == user.START {
				// Получаем все задачи данного пользователя.
				_, reply, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatID, reply))
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Вы не закончили ввод данных."))
			}
			continue
		case "delete_fast_task":
			if userStates[userId].Code == user.START {
				bot.Send(tgbotapi.NewMessage(chatID,
					"Какая из этих задач уже выполнена? (введите её порядковый номер)"))
				_, output, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatID, output))

				user.UpdateUser(userId, userName, user.FAST_TASK_DELETE, "")
				userStates[userId] = user.State{Code: user.FAST_TASK_DELETE, Request: ""}
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Вы не закончили ввод данных."))
			}
			continue
			//case "fill_schedule":
			//	bot.Send(tgbotapi.NewMessage(chatID, "Итак, давайте пробежимся по дням недели "+
			//		"и заполним расписание на каждый из них."))
			//	// Заполнение всех полей расписания в базе данных расписания.
			//	assigneeSchedule, err := schedule.InitScheduleTable(userId) // в основном нам нужны отсюда sch_id для заполнения бд.
			//	if err != nil {
			//		log.Fatal(err)
			//	}
			//	// TODO: Тут должно происходить заполнение полного расписания на неделю.
			//	var scheduleTasks []models.ScheduleTask // sch_id -> Массив заданий на день.
			//	for _, weekdaySch := range assigneeSchedule {
			//		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Заполним расписание на %s\n Введите число дел на этот день.",
			//			services.ParseWeekdayToStr(weekdaySch.WeekDay))))
			//		// <...>
			//	}
			//
			//	err = schedule.FillSchedule(userId, scheduleTasks)
			//	if err != nil {
			//		log.Fatal(err)
			//	}
			//
			//	bot.Send(tgbotapi.NewMessage(chatID, "Здорово! Ваше расписание успешно заполнено! "))
			//}

		}

		if userStates[userId].Code != user.START {
			if userStates[userId].Code == user.FAST_TASK_ENTER_TITLE {
				var fastTask models.FastTask
				if msg.Text == "" {
					bot.Send(tgbotapi.NewMessage(chatID, "Нет текстового сообщения, введите команду заново."))
					user.ResetState(userId, userName, &userStates)
					continue
				}
				fastTask.TaskName = msg.Text
				b, err := json.Marshal(fastTask)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatID,
					"Введите, с какой периодичностью вам будут приходить сообщения. (Например: 1h10m40s)"))

				state := user.State{Code: user.FAST_TASK_ENTER_INTERVAL, Request: string(b)}
				user.UpdateUser(userId, userName, state.Code, state.Request)
				userStates[userId] = state
			} else if userStates[userId].Code == user.FAST_TASK_ENTER_INTERVAL {
				var fastTask models.FastTask
				interval, err := time.ParseDuration(update.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID,
						"Кажется, введённое вами сообщение не удовлетворяет формату. Введите команду ещё раз."))
					user.ResetState(userId, userName, &userStates)
					continue
				}
				currUser, err := user.GetUser(userId)
				if err != nil {
					log.Fatal(err)
				}
				data := []byte(currUser.StateRequest)

				err = json.Unmarshal(data, &fastTask)
				if err != nil {
					log.Fatal(err)
				}
				fastTask.NotifyInterval = interval

				err = fast_task.AddFastTask(userId, fastTask.TaskName, chatID, fastTask.NotifyInterval)

				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatID, "Задача успешно добавлена!"))
				user.ResetState(userId, userName, &userStates)
			} else if userStates[userId].Code == user.FAST_TASK_DELETE {
				fastTasks, _, err := fast_task.OutputFastTasks(userId)

				// Считываем порядковый номер задачи, которую нужно удалить.
				ftNumber, err := strconv.Atoi(msg.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Кажется, вы ввели не число. Введите команду ещё раз."))
					user.ResetState(userId, userName, &userStates)
					continue
				}

				if ftNumber <= 0 || ftNumber > len(fastTasks) {
					bot.Send(tgbotapi.NewMessage(chatID,
						"Кажется, такого дела не существует. Введите команду ещё раз."))
					user.ResetState(userId, userName, &userStates)
					continue
				}

				fastTaskDeleteUrl := DefaultServiceUrl +
					fmt.Sprintf("%v/fast_task/%v", userId, fastTasks[ftNumber-1].Id)

				_, err = utils.Delete(fastTaskDeleteUrl)

				if err != nil {
					log.Fatal(err)
				}

				_, output, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Задача %v успешно удалена!\n", ftNumber)+output))
				user.ResetState(userId, userName, &userStates)
			}
		}
	}
}
