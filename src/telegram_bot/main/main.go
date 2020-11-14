package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"todo_web_service/src/telegram_bot/schedule"
	"todo_web_service/src/telegram_bot/utils"

	"github.com/Syfaro/telegram-bot-api"
	"todo_web_service/src/telegram_bot/fast_task"
	"todo_web_service/src/telegram_bot/user"
)

const (
	SuburbanServiceUrl = utils.DefaultServiceUrl + "suburban"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
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

	// Храним состояния пользователей.
	userStates := make(map[int]user.State)
	err = user.GetStates(&userStates)

	if err != nil {
		log.Fatal(err)
	}
	// читаем обновления из канала
	for update := range updates {
		chatId := update.Message.Chat.ID
		userId := update.Message.From.ID
		userName := update.Message.From.UserName
		userStateCode := userStates[userId].Code

		switch update.Message.Command() {
		case "start":
			if user.IsStartState(userStateCode) {
				err = user.InitUser(userId, userName)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("Здравствуйте, %s.\nДобро пожаловать!", userName)))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "userinfo":
			if user.IsStartState(userStateCode) {
				user, err := user.GetUser(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("Здравствуйте, %s. \nВаш 🆔: %s",
					user.UserName, strconv.Itoa(user.Id))))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "suburban":
			if user.IsStartState(userStateCode) {
				resp, err := http.Get(SuburbanServiceUrl)
				if err != nil {
					log.Fatal(err)
				}

				body, _ := ioutil.ReadAll(resp.Body)

				bot.Send(tgbotapi.NewMessage(chatId, string(body)))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_fast_task":
			if user.IsStartState(userStateCode) {
				state := user.State{Code: user.FAST_TASK_ENTER_TITLE, Request: "{}"}
				user.SetState(userId, userName, &userStates, state)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название нового задания."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "fast_tasks":
			if user.IsStartState(userStateCode) {
				// Получаем все задачи данного пользователя.
				_, reply, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, reply))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_fast_task":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId,
					"Какая из этих задач уже выполнена? (введите её порядковый номер)"))
				_, output, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, output))
				user.SetState(userId, userName, &userStates, user.State{Code: user.FAST_TASK_DELETE_ENTER_NUM, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "fill_schedule":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId, "Выберете день недели, куда вы хотели юы добавить дело:\n"+
					"Понедельник /add_to_mon \nВторник /add_to_tue \nСреда /add_to_wed "+
					"\nЧетверг /add_to_thu \nПятница /add_to_fri \nСуббота /add_to_sat \nВоскресенье /add_to_sun"))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "full_schedule":
			if user.IsStartState(userStateCode) {
				schedule.GetFullSchedule(&bot, userId, chatId)
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
		case "today_schedule":
			if user.IsStartState(userStateCode) {
				_, output, err := schedule.GetWeekdaySchedule(userId, time.Now().Weekday())
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, output))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "tomorrow_schedule":
			if user.IsStartState(userStateCode) {
				_, output, err := schedule.GetWeekdaySchedule(userId, schedule.NextWeekday(time.Now().Weekday()))
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, output))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "weekday_schedule":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId, "На какой день недели вы хотите увидеть расписание?"))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_ENTER_OUTPUT_WEEKDAY, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_schedule":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId,
					"Если вы хотите очистить расписание на всю неделю, используйте /clear_schedule"))
				bot.Send(tgbotapi.NewMessage(chatId,
					"Если вам необходимо очистить расписание на конкретный день недели, используйте \n/delete_weekday_schedule"))
				bot.Send(tgbotapi.NewMessage(chatId,
					"Если вам просто нужно удалить какую-то задачу на конкретный день недели, используйте \n/delete_schedule_task"))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_mon":
			if user.IsStartState(userStateCode) {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_MON)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_tue":
			if user.IsStartState(userStateCode) {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_TUE)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_wed":
			if user.IsStartState(userStateCode) {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_WED)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_thu":
			if user.IsStartState(userStateCode) {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_THU)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_fri":
			if user.IsStartState(userStateCode) {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_FRI)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_sat":
			if user.IsStartState(userStateCode) {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SAT)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_sun":
			if user.IsStartState(userStateCode) {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SUN)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "update_schedule_task":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId, "Введите день недели, в котором нужно обновить задачу."))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_UPDATE_ENTER_WEEKDAY, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_schedule_task":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId, "Введите день недели, в котором нужно удалить задачу."))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_DELETE_WEEKDAY_TASK, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "clear_schedule":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы точно хотите ПОЛНОСТЬЮ очистить ваше текущее расписание? Да или нет?"))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_DELETE_CLEARALL, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_weekday_schedule":
			if user.IsStartState(userStateCode) {
				bot.Send(tgbotapi.NewMessage(chatId, "Расписание на какой день недели вы хотите очистить?"))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_DELETE_WEEKDAY, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "reset":
			user.ResetState(userId, userName, &userStates)
			bot.Send(tgbotapi.NewMessage(chatId, "Ввод данных прерван."))
			continue
		}

		// Если состояние пользователя не начальное.
		if userStateCode != user.START {
			/* FastTask */
			switch userStateCode {
			case user.FAST_TASK_ENTER_TITLE:
				if !fast_task.EnterTitle(&update, &bot, &userStates) {
					continue
				}
			case user.FAST_TASK_ENTER_INTERVAL:
				if !fast_task.EnterInterval(&update, &bot, &userStates) {
					continue
				}
			case user.FAST_TASK_DELETE_ENTER_NUM:
				if !fast_task.EnterDeleteNum(&update, &bot, &userStates) {
					continue
				}
				/* Schedule */
			case user.SCHEDULE_ENTER_TITLE:
				if !schedule.EnterTitle(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_ENTER_PLACE:
				if !schedule.EnterPlace(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_ENTER_SPEAKER:
				if !schedule.EnterSpeaker(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_ENTER_START:
				if !schedule.EnterStart(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_ENTER_END:
				if !schedule.EnterEnd(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_ENTER_OUTPUT_WEEKDAY:
				if !schedule.EnterOutputWeekday(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_UPDATE_ENTER_WEEKDAY:
				if !schedule.EnterUpdateWeekday(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_UPDATE_ENTER_NUM_TASK:
				if !schedule.EnterUpdateNumTask(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_UPDATE_ENTER_TITLE:
				if !schedule.EnterUpdateTitle(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_UPDATE_ENTER_PLACE:
				if !schedule.EnterUpdatePlace(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_UPDATE_ENTER_SPEAKER:
				if !schedule.EnterUpdateSpeaker(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_UPDATE_ENTER_START:
				if !schedule.EnterUpdateStart(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_UPDATE_ENTER_END:
				if !schedule.EnterUpdateEnd(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_DELETE_CLEARALL:
				if !schedule.EnterClearAll(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_DELETE_WEEKDAY_TASK:
				if !schedule.EnterDeleteWeekdayTask(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_DELETE_NUM_TASK:
				if !schedule.EnterDeleteNumTask(&update, &bot, &userStates) {
					continue
				}
			case user.SCHEDULE_DELETE_WEEKDAY:
				if !schedule.EnterDeleteWeekday(&update, &bot, &userStates) {
					continue
				}
			}
		}
	}
}
