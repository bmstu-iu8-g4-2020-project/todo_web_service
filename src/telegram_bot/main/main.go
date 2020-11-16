package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Syfaro/telegram-bot-api"

	"todo_web_service/src/telegram_bot/fast_task"
	"todo_web_service/src/telegram_bot/schedule"
	"todo_web_service/src/telegram_bot/user"
	"todo_web_service/src/telegram_bot/utils"
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
	log.Printf("Authorized on account %s (@%s)", bot.Self.FirstName, bot.Self.UserName)

	// Инициализируем канал, куда будут прилетать обновления от API
	var userConfig = tgbotapi.NewUpdate(0)
	userConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(userConfig)
	if err != nil {
		log.Fatal(err)
	}

	// В отдельном потоке проверяем, прошли ли дедлайны fast_task'ов.
	go fast_task.CheckFastTasks(&bot)

	// Храним состояния пользователей.
	userStates := make(map[int]user.State)
	err = user.GetStates(&userStates)
	if err != nil {
		log.Fatal(err)
	}

	stateFuncDict := make(map[int]user.StateFunc)
	fast_task.FillFastTaskFuncs(&stateFuncDict)
	schedule.FillScheduleFuncs(&stateFuncDict)

	// Читаем обновления из канала
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
				_, _ = bot.Send(tgbotapi.NewMessage(chatId,
					fmt.Sprintf("Здравствуйте, %s.\nДобро пожаловать!", userName)))
				bot.Send(tgbotapi.NewStickerShare(chatId, utils.StickerWelcome))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "userinfo":
			if user.IsStartState(userStateCode) {
				respUser, err := user.GetUser(userId)
				if err != nil {
					log.Fatal(err)
				}
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("%s, \nВаш 🆔: %s",
					respUser.UserName, strconv.Itoa(respUser.Id))))
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

				_, _ = bot.Send(tgbotapi.NewMessage(chatId, string(body)))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_fast_task":
			if user.IsStartState(userStateCode) {
				state := user.State{Code: user.FAST_TASK_ENTER_TITLE, Request: "{}"}
				_ = user.SetState(userId, userName, &userStates, state)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiFastTask+
					"Введите название нового задания."))
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
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, reply))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_fast_task":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiNumber+
					"Какая из этих задач уже выполнена? (введите её порядковый номер)"))
				_, output, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, output))
				_ = user.SetState(userId, userName, &userStates,
					user.State{Code: user.FAST_TASK_DELETE_ENTER_NUM, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "fill_schedule":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiWeekday+
					"Выберете день недели, куда вы хотели бы добавить дело:\n"+
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

				_, _ = bot.Send(tgbotapi.NewMessage(chatId, output))
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

				_, _ = bot.Send(tgbotapi.NewMessage(chatId, output))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "weekday_schedule":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiWeekday+
					"На какой день недели вы хотите увидеть расписание?"))
				_ = user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_ENTER_OUTPUT_WEEKDAY, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_schedule":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId,
					"Если вы хотите очистить расписание на всю неделю, используйте /clear_schedule"))
				_, _ = bot.Send(tgbotapi.NewMessage(chatId,
					"Если вам необходимо очистить расписание на конкретный день недели, используйте "+
						"\n/delete_weekday_schedule"))
				_, _ = bot.Send(tgbotapi.NewMessage(chatId,
					"Если вам просто нужно удалить какую-то задачу на конкретный день недели, используйте "+
						"\n/delete_schedule_task"))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_mon":
			if user.IsStartState(userStateCode) {
				_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_MON)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_tue":
			if user.IsStartState(userStateCode) {
				_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_TUE)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_wed":
			if user.IsStartState(userStateCode) {
				_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_WED)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_thu":
			if user.IsStartState(userStateCode) {
				_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_THU)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_fri":
			if user.IsStartState(userStateCode) {
				_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_FRI)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_sat":
			if user.IsStartState(userStateCode) {
				_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SAT)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "add_to_sun":
			if user.IsStartState(userStateCode) {
				_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SUN)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "update_schedule_task":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiWeekday+
					"Введите день недели, в котором нужно обновить задачу."))
				_ = user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_UPDATE_ENTER_WEEKDAY, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_schedule_task":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiWeekday+
					"Введите день недели, в котором нужно удалить задачу."))
				_ = user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_DELETE_WEEKDAY_TASK, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "clear_schedule":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId,
					"Вы точно хотите ПОЛНОСТЬЮ очистить ваше текущее расписание? Да или нет?"))
				_ = user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_DELETE_CLEARALL, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "delete_weekday_schedule":
			if user.IsStartState(userStateCode) {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiWeekday+
					"Расписание на какой день недели вы хотите очистить?"))
				_ = user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_DELETE_WEEKDAY, Request: "{}"})
			} else {
				user.SendEnteringNotFinished(&bot, chatId)
			}
			continue
		case "reset":
			if !user.IsStartState(userStateCode) {
				_ = user.ResetState(userId, userName, &userStates)
				_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiReset+"Ввод данных прерван."))
			} else {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId,
					utils.EmojiWarning+"Вы не вводите данные. Вам нечего прерывать"))
			}
			continue
		}

		// Если состояние пользователя не начальное.
		if !user.IsStartState(userStateCode) {
			stateFuncDict[userStateCode](&update, &bot, &userStates)
		}
	}
}
