package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Syfaro/telegram-bot-api"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/fast_task"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/schedule"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/user"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/utils"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/weather"
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
	weather.FillWeatherFuncs(&stateFuncDict)

	var chatId int64
	var userId int
	var userName string

	// Читаем обновления из канала
	for update := range updates {
		if update.CallbackQuery != nil {
			chatId = update.CallbackQuery.Message.Chat.ID
			userId = update.CallbackQuery.From.ID
			userName = update.CallbackQuery.From.UserName
		} else {
			chatId = update.Message.Chat.ID
			userId = update.Message.From.ID
			userName = update.Message.From.UserName
		}

		userStateCode := userStates[userId].Code

		if update.Message != nil {
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
			//case "suburban":
			//	if user.IsStartState(userStateCode) {
			//		resp, err := http.Get(SuburbanServiceUrl)
			//		if err != nil {
			//			log.Fatal(err)
			//		}
			//
			//		body, _ := ioutil.ReadAll(resp.Body)
			//
			//		_, _ = bot.Send(tgbotapi.NewMessage(chatId, string(body)))
			//	} else {
			//		user.SendEnteringNotFinished(&bot, chatId)
			//	}
			//	continue
			case "weather":
				if user.IsStartState(userStateCode) {
					msg := tgbotapi.NewMessage(chatId, "Я могу предоставить вам данные о погоде:")
					msg.ReplyMarkup = WeatherKeyboard
					_, _ = bot.Send(msg)
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
					msg := tgbotapi.NewMessage(chatId,
						utils.EmojiWeekday+"Выберете день недели, куда вы хотели бы добавить дело:")
					msg.ReplyMarkup = WeekdayScheduleKeyboard
					_, _ = bot.Send(msg)
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
					msg := tgbotapi.NewMessage(chatId, "Выберите опцию удаления:")
					msg.ReplyMarkup = ScheduleDeleteKeyboard
					_, _ = bot.Send(msg)
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
		}

		if update.CallbackQuery != nil {
			if user.IsStartState(userStateCode) {
				switch update.CallbackQuery.Data {
				case "add_mon":
					_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_MON)
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
				case "add_tue":
					_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_TUE)
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
				case "add_wed":
					_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_WED)
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
				case "add_thu":
					_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_THU)
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
				case "add_fri":
					_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_FRI)
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
				case "add_sat":
					_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SAT)
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
				case "add_sun":
					_ = schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SUN)
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiTitle+"Введите название дела."))
				case "delete_schedule_task":
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiWeekday+
						"Введите день недели, в котором нужно удалить задачу."))
					_ = user.SetState(userId, userName, &userStates,
						user.State{Code: user.SCHEDULE_DELETE_WEEKDAY_TASK, Request: "{}"})
				case "clear_weekday_schedule":
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiWeekday+
						"Расписание на какой день недели вы хотите очистить?"))
					_ = user.SetState(userId, userName, &userStates,
						user.State{Code: user.SCHEDULE_DELETE_WEEKDAY, Request: "{}"})
				case "clear_schedule":
					_, _ = bot.Send(tgbotapi.NewMessage(chatId,
						fmt.Sprintf("Вы точно хотите %sПОЛНОСТЬЮ%s очистить ваше текущее расписание? Да или нет?",
							utils.EmojiWarning, utils.EmojiWarning)))
					_ = user.SetState(userId, userName, &userStates,
						user.State{Code: user.SCHEDULE_DELETE_CLEARALL, Request: "{}"})
				case "current_weather":
					msg := tgbotapi.NewMessage(chatId, "Как бы вы хотели получить данные о погоде?")
					msg.ReplyMarkup = WeatherChooseCurrentKeyboard
					_, _ = bot.Send(msg)
				case "weather_forecast":
					msg := tgbotapi.NewMessage(chatId, "Как бы вы хотели получить прогноз погоды?")
					msg.ReplyMarkup = WeatherChooseForecastKeyboard
					_, _ = bot.Send(msg)
				case "curr_location":
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf(
						"Пришлите мне свою геопозицию. \n(нажмите на %s и выберите \"Геопозиция\")",
						utils.EmojiPaperclip)))
					_ = user.SetState(userId, userName, &userStates,
						user.State{Code: user.WEATHER_CURRENT_SEND_LOCATION, Request: "{}"})
				case "curr_place_name":
					_, _ = bot.Send(tgbotapi.NewMessage(chatId, utils.EmojiLocation+
						"Введите место, где вы бы хотели узнать данные о погоде."))
					_ = user.SetState(userId, userName, &userStates,
						user.State{Code: user.WEATHER_CURRENT_SEND_NAME, Request: "{}"})
				case "forecast_location":
					_, _ = (*bot).Send(tgbotapi.NewMessage(chatId, fmt.Sprintf(
						"Пришлите мне свою геопозицию. \n(нажмите на %s и выберите \"Геопозиция\")",
						utils.EmojiPaperclip)))
					_ = user.SetState(userId, userName, &userStates,
						user.State{Code: user.WEATHER_FORECAST_SEND_LOCATION, Request: "{}"})
				case "forecast_place_name":
					_, _ = (*bot).Send(tgbotapi.NewMessage(chatId, utils.EmojiLocation+
						"Введите место, где вы бы хотели узнать данные о погоде."))
					_ = user.SetState(userId, userName, &userStates,
						user.State{Code: user.WEATHER_FORECAST_SEND_NAME, Request: "{}"})
				}
				continue
			} else {
				_, _ = bot.Send(tgbotapi.NewMessage(chatId,
					utils.EmojiWarning+"Вы пытаетесь использовать кнопки во время ввода данных. "+
						"Закончите ввод, либо используйте /reset чтобы его прервать."))
			}
		}

		// Если состояние пользователя не начальное, используем функцию текущего состояния.
		if !user.IsStartState(userStateCode) {
			stateFuncDict[userStateCode](&update, &bot, &userStates)
		}
	}
}


var WeekdayScheduleKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Понедельник", "add_mon"),
		tgbotapi.NewInlineKeyboardButtonData("Четверг", "add_thu"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вторник", "add_tue"),
		tgbotapi.NewInlineKeyboardButtonData("Пятница", "add_fri"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Среда", "add_wed"),
		tgbotapi.NewInlineKeyboardButtonData("Суббота", "add_sat"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Воскресенье", "add_sun"),
	),
)

var ScheduleDeleteKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiTitle+
			"Удалить задачу", "delete_schedule_task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiWeekday+
			"Очистить расписание на день", "clear_weekday_schedule"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiFire+
			"Полностью очистить расписание", "clear_schedule"),
	),
)

var WeatherKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("На текущий момент времени", "current_weather"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Прогноз на ближайшие 5 дней", "weather_forecast"),
	),
)

var WeatherChooseCurrentKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiLocation+"По геопозиции", "curr_location"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiMap+"По введённому месту", "curr_place_name"),
	),
)

var WeatherChooseForecastKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiLocation+"По геопозиции", "forecast_location"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(utils.EmojiMap+"По введённому месту", "forecast_place_name"),
	),
)

