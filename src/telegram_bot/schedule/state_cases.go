package schedule

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"

	"todo_web_service/src/models"
	"todo_web_service/src/services"
	"todo_web_service/src/telegram_bot/user"
	"todo_web_service/src/telegram_bot/utils"
)

func FillScheduleFuncs(stateFuncDict *map[int]user.StateFunc) {
	(*stateFuncDict)[user.SCHEDULE_ENTER_TITLE] = EnterTitle
	(*stateFuncDict)[user.SCHEDULE_ENTER_PLACE] = EnterPlace
	(*stateFuncDict)[user.SCHEDULE_ENTER_SPEAKER] = EnterSpeaker
	(*stateFuncDict)[user.SCHEDULE_ENTER_START] = EnterStart
	(*stateFuncDict)[user.SCHEDULE_ENTER_END] = EnterEnd
	(*stateFuncDict)[user.SCHEDULE_ENTER_OUTPUT_WEEKDAY] = EnterOutputWeekday
	(*stateFuncDict)[user.SCHEDULE_UPDATE_ENTER_WEEKDAY] = EnterUpdateWeekday
	(*stateFuncDict)[user.SCHEDULE_UPDATE_ENTER_NUM_TASK] = EnterUpdateNumTask
	(*stateFuncDict)[user.SCHEDULE_UPDATE_ENTER_TITLE] = EnterUpdateTitle
	(*stateFuncDict)[user.SCHEDULE_UPDATE_ENTER_PLACE] = EnterUpdatePlace
	(*stateFuncDict)[user.SCHEDULE_UPDATE_ENTER_SPEAKER] = EnterUpdateSpeaker
	(*stateFuncDict)[user.SCHEDULE_UPDATE_ENTER_START] = EnterUpdateStart
	(*stateFuncDict)[user.SCHEDULE_UPDATE_ENTER_END] = EnterUpdateEnd
	(*stateFuncDict)[user.SCHEDULE_DELETE_CLEARALL] = EnterClearAll
	(*stateFuncDict)[user.SCHEDULE_DELETE_WEEKDAY_TASK] = EnterDeleteWeekdayTask
	(*stateFuncDict)[user.SCHEDULE_DELETE_NUM_TASK] = EnterDeleteNumTask
	(*stateFuncDict)[user.SCHEDULE_DELETE_WEEKDAY] = EnterDeleteWeekday
}

func EnterTitle(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	if update.Message.Text == "" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы отправили не текстовое сообщение. Введите название задания."))
		return
	}

	scheduleTask.Title = update.Message.Text
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите место проведения."))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_ENTER_PLACE, Request: string(b)})
}

func EnterPlace(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	if update.Message.Text == "" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы отправили не текстовое сообщение. Введите место проведения."))
		return
	}

	scheduleTask.Place = update.Message.Text
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите имя спикера. (преподавателя, лектора, выступающего)"))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_ENTER_SPEAKER, Request: string(b)})
}

func EnterSpeaker(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	if update.Message.Text == "" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы отправили не текстовое сообщение. Введите имя спикера."))
		return
	}

	scheduleTask.Speaker = update.Message.Text
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите время начала дела. (например: 10:00)\n"))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_ENTER_START, Request: string(b)})
}

func EnterStart(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	startTime, err := time.Parse(LayoutTime, update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ой, кажется, вы ввели время не в подходящем формате. "+
			"Попробуйте ещё раз"))
		return
	}
	scheduleTask.Start = startTime
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите время окончания дела. (например: 19:00)"))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_ENTER_END, Request: string(b)})
}

func EnterEnd(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask

	err := json.Unmarshal([]byte((*userStates)[update.Message.From.ID].Request), &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	endTime, err := time.Parse(LayoutTime, update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ой, кажется, вы ввели время не в подходящем формате. "+
			"Попробуйте ещё раз"))
		return
	}

	// Время конца дела должно быть после времени начала.
	if !endTime.After(scheduleTask.Start) {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Время окончания дела не может быть раньше времени его начала. Попробуйте ещё раз."))
		return
	}

	scheduleTask.End = endTime

	err = AddScheduleTask(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Супер! %s пополнился новой задачей.",
		services.WeekdayToStr(scheduleTask.WeekDay))))

	user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
}

func EnterOutputWeekday(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	weekday, err := services.StrToWeekday(strings.Title(update.Message.Text))
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Нет-нет. Введите день недели. (например: Понедельник)"))
		return
	}

	_, output, err := GetWeekdaySchedule(update.Message.From.ID, weekday)
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, output))

	user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
}

func EnterClearAll(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	if strings.ToLower(update.Message.Text) == "да" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ок, очищаю ваше расписание..."))
		err := ClearAll(update.Message.From.ID)
		if err != nil {
			log.Fatal(err)
		}

		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Расписание очищено!"))
		user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
	} else if strings.ToLower(update.Message.Text) == "нет" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Хорошо, не будем ничего удалять."))
		user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
	} else {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ответ не понятен, введите да, либо нет."))
	}
}

func EnterDeleteWeekdayTask(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	weekday, err := services.StrToWeekday(strings.Title(update.Message.Text))
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Нет-нет. Введите день недели. (например: Понедельник)"))
		return
	}

	weekdaySchedule, output, err := GetWeekdaySchedule(update.Message.From.ID, weekday)
	if err != nil {
		log.Fatal(err)
	}
	if weekdaySchedule == nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Кажется, на %s задач не существует. Удалять тут нечего. Ещё разок? /clear_schedule_task",
				strings.ToLower(update.Message.Text))))
		user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)

		return
	}
	b, err := json.Marshal(weekdaySchedule)
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, output))
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Итак, теперь введите номер задачи, которую вы желаете удалить."))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates, user.State{Code: user.SCHEDULE_DELETE_NUM_TASK, Request: string(b)})
}

func EnterDeleteNumTask(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTasks []models.ScheduleTask

	err := json.Unmarshal([]byte((*userStates)[update.Message.From.ID].Request), &scheduleTasks)
	if err != nil {
		log.Fatal(err)
	}
	weekday := scheduleTasks[0].WeekDay

	num, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы ввели не число. Введите номер задания, который хотите удалить."))
		return
	}

	if num <= 0 || num > len(scheduleTasks) {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, такого дела не существует. Введите номер задания, который хотите удалить."))
		return
	}

	_, err = utils.Delete(fmt.Sprintf("%s%v/schedule/%v/", utils.DefaultServiceUrl, update.Message.From.ID, scheduleTasks[num-1].Id))
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задание успешно удалено."))
	_, output, err := GetWeekdaySchedule(update.Message.From.ID, weekday)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, output))

	user.ResetState(update.Message.From.ID, update.Message.Chat.UserName, userStates)
}

func EnterDeleteWeekday(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	weekday, err := services.StrToWeekday(strings.Title(update.Message.Text))
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Нет-нет. Введите день недели. (например: Понедельник)"))
		return
	}

	scheduleTasks, _, err := GetWeekdaySchedule(update.Message.From.ID, weekday)
	if err != nil {
		log.Fatal(err)
	}
	if scheduleTasks == nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Кажется, текущий день недели итак пуст.\n"+
			"Введите другой день недели, либо прервите ввод -- /reset"))
		return
	}

	_, err = utils.Delete(fmt.Sprintf("%s%v/schedule/delete/%s/",
		utils.DefaultServiceUrl, update.Message.From.ID, weekday.String()))
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s больше не имеет задач. Они успешно очищены.",
		services.WeekdayToStr(weekday))))

	user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
}

func EnterUpdateWeekday(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	weekday, err := services.StrToWeekday(strings.Title(update.Message.Text))
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Нет-нет. Введите день недели. (например: Понедельник)"))
		return
	}

	weekdaySchedule, output, err := GetWeekdaySchedule(update.Message.From.ID, weekday)
	if err != nil {
		log.Fatal(err)
	}
	if weekdaySchedule == nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Кажется, на %s задач не существует. Обновлять тут нечего. Ещё разок? /update_schedule_task",
				strings.ToLower(update.Message.Text))))
		user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
	}
	b, err := json.Marshal(weekdaySchedule)
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, output))
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Итак, теперь введите номер задачи, которую вы желаете обновить."))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates, user.State{Code: user.SCHEDULE_UPDATE_ENTER_NUM_TASK, Request: string(b)})
}

func EnterUpdateNumTask(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTasks []models.ScheduleTask
	err := json.Unmarshal([]byte((*userStates)[update.Message.From.ID].Request), &scheduleTasks)
	if err != nil {
		log.Fatal(err)
	}

	num, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы ввели не число. Введите номер задания, который хотите удалить."))
		return
	}

	if num <= 0 || num > len(scheduleTasks) {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, такого дела не существует. Введите номер задания, который хотите удалить."))
		return
	}
	b, err := json.Marshal(scheduleTasks[num-1])
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ок. Введите новое название дела."))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates, user.State{Code: user.SCHEDULE_UPDATE_ENTER_TITLE, Request: string(b)})
}

func EnterUpdateTitle(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	if update.Message.Text == "" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы отправили не текстовое сообщение. Введите название задания."))
		return
	}

	scheduleTask.Title = update.Message.Text
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите новое место проведения."))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_UPDATE_ENTER_PLACE, Request: string(b)})
}

func EnterUpdatePlace(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	if update.Message.Text == "" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы отправили не текстовое сообщение. Введите место проведения."))
		return
	}

	scheduleTask.Place = update.Message.Text
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите имя спикера. (преподавателя, лектора, выступающего)"))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_UPDATE_ENTER_SPEAKER, Request: string(b)})
}

func EnterUpdateSpeaker(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	if update.Message.Text == "" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы отправили не текстовое сообщение. Введите имя спикера."))
		return
	}

	scheduleTask.Speaker = update.Message.Text
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите новое время начала дела. (например: 10:00)\n"))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_UPDATE_ENTER_START, Request: string(b)})
}

func EnterUpdateStart(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask
	data := []byte((*userStates)[update.Message.From.ID].Request)

	err := json.Unmarshal(data, &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	startTime, err := time.Parse(LayoutTime, update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Ой, кажется, вы ввели время не в подходящем формате. Попробуйте ещё раз"))
		return
	}
	scheduleTask.Start = startTime
	b, err := json.Marshal(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите новое время окончания дела. (например: 19:00)"))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.SCHEDULE_UPDATE_ENTER_END, Request: string(b)})
}

func EnterUpdateEnd(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) {
	var scheduleTask models.ScheduleTask

	err := json.Unmarshal([]byte((*userStates)[update.Message.From.ID].Request), &scheduleTask)
	if err != nil {
		log.Fatal(err)
	}
	endTime, err := time.Parse(LayoutTime, update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ой, кажется, вы ввели время не в подходящем формате. "+
			"Попробуйте ещё раз"))
		return
	}

	// Время конца дела должно быть после времени начала.
	if !endTime.After(scheduleTask.Start) {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Время окончания дела не может быть раньше времени его начала. Попробуйте ещё раз."))
		return
	}

	scheduleTask.End = endTime

	err = UpdateScheduleTask(scheduleTask)
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s: задание обновлено успешно.",
		services.WeekdayToStr(scheduleTask.WeekDay))))

	user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)
}
