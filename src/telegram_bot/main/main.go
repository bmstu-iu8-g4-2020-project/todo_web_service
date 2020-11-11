package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"todo_web_service/src/models"
	"todo_web_service/src/services"
	"todo_web_service/src/telegram_bot/schedule"
	"todo_web_service/src/telegram_bot/utils"

	"github.com/Syfaro/telegram-bot-api"
	"todo_web_service/src/telegram_bot/fast_task"
	"todo_web_service/src/telegram_bot/user"
)

const (
	DefaultServiceUrl  = "http://localhost:8080/"
	SuburbanServiceUrl = DefaultServiceUrl + "suburban"
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

	// Храним состояния пользователей в бд.
	userStates := make(map[int]user.State)
	err = user.GetStates(&userStates)
	if err != nil {
		log.Fatal(err)
	}
	// читаем обновления из канала
	for update := range updates {
		chatId := update.Message.Chat.ID
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
				bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("Здравствуйте, %s.\nДобро пожаловать!", userName)))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "userinfo":
			if userStates[userId].Code == user.START {
				user, err := user.GetUser(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("Здравствуйте, %s. \nВаш 🆔: %s",
					user.UserName, strconv.Itoa(user.Id))))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "suburban":
			if userStates[userId].Code == user.START {
				resp, err := http.Get(SuburbanServiceUrl)
				if err != nil {
					log.Fatal(err)
				}

				body, _ := ioutil.ReadAll(resp.Body)

				bot.Send(tgbotapi.NewMessage(chatId, string(body)))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_fast_task":
			if userStates[userId].Code == user.START {
				state := user.State{Code: user.FAST_TASK_ENTER_TITLE, Request: "{}"}
				user.SetState(userId, userName, &userStates, state)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название нового задания."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "fast_tasks":
			if userStates[userId].Code == user.START {
				// Получаем все задачи данного пользователя.
				_, reply, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, reply))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "delete_fast_task":
			if userStates[userId].Code == user.START {
				bot.Send(tgbotapi.NewMessage(chatId,
					"Какая из этих задач уже выполнена? (введите её порядковый номер)"))
				_, output, err := fast_task.OutputFastTasks(userId)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, output))
				user.SetState(userId, userName, &userStates, user.State{Code: user.FAST_TASK_DELETE, Request: "{}"})
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "fill_schedule":
			if userStates[userId].Code == user.START {
				bot.Send(tgbotapi.NewMessage(chatId, "Выберете день недели, куда вы хотели юы добавить дело:\n"+
					"Понедельник /add_to_mon \nВторник /add_to_tue \nСреда /add_to_wed "+
					"\nЧетверг /add_to_thu \nПятница /add_to_fri \nСуббота /add_to_sat \nВоскресенье /add_to_sun"))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "today_schedule":
			if userStates[userId].Code == user.START {
				_, output, err := schedule.GetWeekdaySchedule(userId, time.Now().Weekday())
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, output))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "tomorrow_schedule":
			if userStates[userId].Code == user.START {
				_, output, err := schedule.GetWeekdaySchedule(userId, time.Now().Weekday()+1)
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, output))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "weekday_schedule":
			if userStates[userId].Code == user.START {
				bot.Send(tgbotapi.NewMessage(chatId, "На какой день недели вы хотите увидеть расписание?"))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_ENTER_OUTPUT_WEEKDAY, Request: "{}"})
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "delete_schedule":
			if userStates[userId].Code == user.START {
				bot.Send(tgbotapi.NewMessage(chatId,
					"Если вы хотите очистить расписание на всю неделю, используйте /clear_schedule"))
				bot.Send(tgbotapi.NewMessage(chatId,
					"Если вам необходимо очистить расписание на конкретный день недели, используйте \n/clear_weekday_schedule (в разработке)"))
				bot.Send(tgbotapi.NewMessage(chatId,
					"Если вам просто нужно удалить какую-то задачу на конкретный день недели, используйте \n/clear_schedule_task (в разработке)"))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_to_mon":
			if userStates[userId].Code == user.START {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_MON)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_to_tue":
			if userStates[userId].Code == user.START {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_TUE)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_to_wed":
			if userStates[userId].Code == user.START {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_WED)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_to_thu":
			if userStates[userId].Code == user.START {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_THU)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_to_fri":
			if userStates[userId].Code == user.START {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_FRI)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_to_sat":
			if userStates[userId].Code == user.START {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SAT)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "add_to_sun":
			if userStates[userId].Code == user.START {
				schedule.AddToWeekday(userId, userName, &userStates, user.SCHEDULE_FILL_SUN)
				bot.Send(tgbotapi.NewMessage(chatId, "Введите название дела."))
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "clear_schedule_task":
			if userStates[userId].Code == user.START {
				bot.Send(tgbotapi.NewMessage(chatId, "Введите день недели, в котором нужно удалить задачу."))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_DELETE_TASK, Request: "{}"})
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			} // TODO: Доделать сие.
			continue

		case "clear_schedule":
			if userStates[userId].Code == user.START {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы точно хотите ПОЛНОСТЬЮ очистить ваше текущее расписание? Да или нет?"))
				user.SetState(userId, userName, &userStates, user.State{Code: user.SCHEDULE_DELETE_CLEARALL, Request: "{}"})
			} else {
				bot.Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
					"Если хотите прервать ввод, используйте /reset."))
			}
			continue
		case "reset":
			user.ResetState(userId, userName, &userStates)
			bot.Send(tgbotapi.NewMessage(chatId, "Ввод данных прерван."))
		}

		if userStates[userId].Code != user.START {

			/* FastTask */
			if userStates[userId].Code == user.FAST_TASK_ENTER_TITLE {
				var fastTask models.FastTask
				if msg.Text == "" {
					bot.Send(tgbotapi.NewMessage(chatId, "Нет текстового сообщения, попробуйте ещё раз."))
					continue
				}
				fastTask.TaskName = msg.Text
				b, err := json.Marshal(fastTask)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId,
					"Введите, с какой периодичностью вам будут приходить сообщения. (Например: 1h10m40s)"))

				user.SetState(userId, userName, &userStates,
					user.State{Code: user.FAST_TASK_ENTER_INTERVAL, Request: string(b)})
			} else if userStates[userId].Code == user.FAST_TASK_ENTER_INTERVAL {
				var fastTask models.FastTask
				interval, err := time.ParseDuration(update.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatId,
						"Кажется, введённое вами сообщение не удовлетворяет формату. (пример: 1h40m13s) Попробуйте ещё раз."))
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

				err = fast_task.AddFastTask(userId, fastTask.TaskName, chatId, fastTask.NotifyInterval)

				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, "Задача успешно добавлена!"))
				user.ResetState(userId, userName, &userStates)
			} else if userStates[userId].Code == user.FAST_TASK_DELETE {
				fastTasks, _, err := fast_task.OutputFastTasks(userId)

				// Считываем порядковый номер задачи, которую нужно удалить.
				ftNumber, err := strconv.Atoi(msg.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatId,
						"Кажется, вы ввели не число. Введите номер задания, который хотите удалить."))
					continue
				}

				if ftNumber <= 0 || ftNumber > len(fastTasks) {
					bot.Send(tgbotapi.NewMessage(chatId,
						"Кажется, такого дела не существует. Введите номер задания, который хотите удалить."))
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

				bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("Задача %v успешно удалена!\n", ftNumber)+output))
				user.ResetState(userId, userName, &userStates)

				/* Schedule */
			} else if userStates[userId].Code == user.SCHEDULE_ENTER_TITLE {
				var scheduleTask models.ScheduleTask
				data := []byte(userStates[userId].Request)

				err = json.Unmarshal(data, &scheduleTask)
				if err != nil {
					log.Fatal(err)
				}

				if msg.Text == "" {
					bot.Send(tgbotapi.NewMessage(chatId,
						"Кажется, вы отправили не текстовое сообщение. Введите название задания."))

					continue
				}

				scheduleTask.Title = msg.Text
				b, err := json.Marshal(scheduleTask)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, "Введите место проведения."))

				user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_ENTER_PLACE, Request: string(b)})
			} else if userStates[userId].Code == user.SCHEDULE_ENTER_PLACE {
				var scheduleTask models.ScheduleTask
				data := []byte(userStates[userId].Request)

				err = json.Unmarshal(data, &scheduleTask)
				if err != nil {
					log.Fatal(err)
				}

				if msg.Text == "" {
					bot.Send(tgbotapi.NewMessage(chatId,
						"Кажется, вы отправили не текстовое сообщение. Введите место проведения."))
					continue
				}

				scheduleTask.Place = msg.Text
				b, err := json.Marshal(scheduleTask)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, "Введите имя спикера. (преподавателя, лектора, выступающего)"))

				user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_ENTER_SPEAKER, Request: string(b)})
			} else if userStates[userId].Code == user.SCHEDULE_ENTER_SPEAKER {
				var scheduleTask models.ScheduleTask
				data := []byte(userStates[userId].Request)

				err = json.Unmarshal(data, &scheduleTask)
				if err != nil {
					log.Fatal(err)
				}

				if msg.Text == "" {
					bot.Send(tgbotapi.NewMessage(chatId,
						"Кажется, вы отправили не текстовое сообщение. Введите имя спикера."))
					continue
				}

				scheduleTask.Speaker = msg.Text
				b, err := json.Marshal(scheduleTask)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, "Введите время начала дела. Например: 10:00AM\n"))

				user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_ENTER_START, Request: string(b)})
			} else if userStates[userId].Code == user.SCHEDULE_ENTER_START {
				var scheduleTask models.ScheduleTask
				data := []byte(userStates[userId].Request)

				err = json.Unmarshal(data, &scheduleTask)
				if err != nil {
					log.Fatal(err)
				}
				startTime, err := time.Parse(time.Kitchen, msg.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatId, "Ой, кажется, вы ввели время не в подходящем формате. "+
						"Попробуйте ещё раз"))
					continue
				}
				scheduleTask.Start = startTime
				b, err := json.Marshal(scheduleTask)
				if err != nil {
					log.Fatal(err)
				}
				bot.Send(tgbotapi.NewMessage(chatId, "Введите время окончания дела. (например: 3:45PM)"))

				user.SetState(userId, userName, &userStates,
					user.State{Code: user.SCHEDULE_ENTER_END, Request: string(b)})
			} else if userStates[userId].Code == user.SCHEDULE_ENTER_END {
				var scheduleTask models.ScheduleTask
				data := []byte(userStates[userId].Request)

				err = json.Unmarshal(data, &scheduleTask)
				if err != nil {
					log.Fatal(err)
				}
				endTime, err := time.Parse(time.Kitchen, msg.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatId, "Ой, кажется, вы ввели время не в подходящем формате. "+
						"Попробуйте ещё раз"))
					continue
				}
				scheduleTask.End = endTime

				err = schedule.AddScheduleTask(scheduleTask)
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("Супер! %s пополнился новой задачей.",
					services.WeekdayToStr(scheduleTask.WeekDay))))

				user.ResetState(userId, userName, &userStates)
			} else if userStates[userId].Code == user.SCHEDULE_ENTER_OUTPUT_WEEKDAY {
				weekday, err := services.StrToWeekday(strings.Title(msg.Text))
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatId, "Нет-нет. Введите день недели. (например: Понедельник)"))
					continue
				}

				_, output, err := schedule.GetWeekdaySchedule(userId, weekday)
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, output))

				user.ResetState(userId, userName, &userStates)
			} else if userStates[userId].Code == user.SCHEDULE_DELETE_CLEARALL {
				if strings.ToLower(msg.Text) == "да" {
					bot.Send(tgbotapi.NewMessage(chatId, "Ок, очищаю ваше расписание..."))
					err := schedule.ClearAll(userId)
					if err != nil {
						log.Fatal(err)
					}

					bot.Send(tgbotapi.NewMessage(chatId, "Расписание очищено!"))
					user.ResetState(userId, userName, &userStates)
				} else if strings.ToLower(msg.Text) == "нет" {
					bot.Send(tgbotapi.NewMessage(chatId, "Хорошо, не будем ничего удалять."))
					user.ResetState(userId, userName, &userStates)
				} else {
					bot.Send(tgbotapi.NewMessage(chatId, "Ответ не понятен, введите да, либо нет."))
				}
				continue
			} else if userStates[userId].Code == user.SCHEDULE_DELETE_TASK {
				weekday, err := services.StrToWeekday(strings.Title(msg.Text))
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatId, "Нет-нет. Введите день недели. (например: Понедельник)"))
					continue
				}

				weekdaySchedule, output, err := schedule.GetWeekdaySchedule(userId, weekday)
				if err != nil {
					log.Fatal(err)
				}
				if weekdaySchedule == nil {
					bot.Send(tgbotapi.NewMessage(chatId,
						fmt.Sprintf("Кажется, на %s задач не существует. Удалять тут нечего. Ещё разок? /clear_schedule_task",
							strings.ToLower(msg.Text))))
				}
				b, err = json.Marshal(weekdaySchedule)
				// TODO: замаршелить weekdaySchedule и отправить его дальше.
				if err != nil {
					log.Fatal(err)
				}

				bot.Send(tgbotapi.NewMessage(chatId, output))
				bot.Send(tgbotapi.NewMessage(chatId, "Итак, теперь введите номер задачи, которую вы желаете удалить из"))

				user.ResetState(userId, userName, &userStates)
			}
		}
	}
}
