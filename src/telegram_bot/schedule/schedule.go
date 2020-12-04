// Copyright 2020 aaaaaaaalesha <sks2311211@mail.ru>

package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/models"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/services"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/user"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/utils"
)

const (
	LayoutTime = "15:04"
)

func AddToWeekday(userId int, userName string, userStates *map[int]user.State, stateCode int) error {
	weekday := services.StateCodeToWeekDay(stateCode)
	b, err := json.Marshal(models.ScheduleTask{AssigneeId: userId, WeekDay: weekday})
	if err != nil {
		return err
	}

	_ = user.SetState(userId, userName, userStates, user.State{Code: user.SCHEDULE_ENTER_TITLE, Request: string(b)})
	return nil
}

func NextWeekday(weekday time.Weekday) time.Weekday {
	if weekday == time.Saturday {
		return time.Sunday
	}

	return weekday + 1
}

func AddScheduleTask(scheduleTask models.ScheduleTask) error {
	bytesRepr, err := json.Marshal(scheduleTask)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%v/schedule/", utils.DefaultServiceUrl, scheduleTask.AssigneeId)

	_, err = http.Post(url, "application/json", bytes.NewBuffer(bytesRepr))

	if err != nil {
		return err
	}

	return nil
}

func UpdateScheduleTask(scheduleTask models.ScheduleTask) error {
	bytesRepr, err := json.Marshal(scheduleTask)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%v/schedule/", utils.DefaultServiceUrl, scheduleTask.AssigneeId)

	_, err = utils.Put(url, bytes.NewBuffer(bytesRepr))

	if err != nil {
		return err
	}

	return nil
}

func GetFullSchedule(bot **tgbotapi.BotAPI, userId int, chatId int64) {
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday}
	for _, weekday := range weekdays {
		_, output, err := GetWeekdaySchedule(userId, weekday)
		if err != nil {
			log.Fatal(err)
		}
		(*bot).Send(tgbotapi.NewMessage(chatId, output))
	}
}

func GetWeekdaySchedule(userId int, weekday time.Weekday) ([]models.ScheduleTask, string, error) {
	url := fmt.Sprintf("%s%v/schedule/%s/", utils.DefaultServiceUrl, userId, weekday.String())

	resp, err := http.Get(url)

	if err != nil {
		return nil, "", err
	}
	var weekdaySchedule []models.ScheduleTask
	err = json.NewDecoder(resp.Body).Decode(&weekdaySchedule)
	if err != nil {
		return nil, "", err
	}

	if weekdaySchedule == nil {
		output := fmt.Sprintf("Похоже, что %s у вас ещё не имеет дел, добавим? \n/fill_schedule",
			strings.ToLower(services.WeekdayToStr(weekday)))
		return nil, output, nil
	}

	// Сортируем по времени начала дел.
	sort.SliceStable(weekdaySchedule, func(i, j int) bool {
		return weekdaySchedule[i].Start.Before(weekdaySchedule[j].Start)
	})

	var output strings.Builder
	_, _ = fmt.Fprintf(&output, "%s %s:\n\n", utils.EmojiWeekday, services.WeekdayToStr(weekday))
	for i, scheduleTask := range weekdaySchedule {
		_, _ = fmt.Fprintf(&output, "Задача %v\n", i+1)
		_, _ = fmt.Fprintf(&output, "%s %s\n", utils.EmojiTitle, scheduleTask.Title)
		_, _ = fmt.Fprintf(&output, "%s %s - %s \n", utils.EmojiTime, scheduleTask.Start.Format(LayoutTime), scheduleTask.End.Format(LayoutTime))
		_, _ = fmt.Fprintf(&output, "%s %s\n", utils.EmojiPlace, scheduleTask.Place)
		_, _ = fmt.Fprintf(&output, "%s %s\n\n", utils.EmojiSpeaker, scheduleTask.Speaker)
	}

	return weekdaySchedule, output.String(), nil
}

func ClearAll(userId int) error {
	_, err := utils.Delete(fmt.Sprintf("%s%v/schedule/", utils.DefaultServiceUrl, userId))
	if err != nil {
		return err
	}
	return nil
}
