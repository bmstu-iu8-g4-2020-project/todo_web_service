package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"todo_web_service/src/models"
	"todo_web_service/src/services"
	"todo_web_service/src/telegram_bot/user"
	"todo_web_service/src/telegram_bot/utils"
)

const (
	DefaultServiceUrl = "http://localhost:8080/"
	emojiTitle        = "📃"
	emojiSpeaker      = "👤"
	emojiPlace        = "🏫"
	emojiTime         = "⌚"
	layoutTime        = "15:04"
)

func AddToWeekday(userId int, userName string, userStates *map[int]user.State, stateCode int) error {
	weekday := services.StateCodeToWeekDay(stateCode)
	b, err := json.Marshal(models.ScheduleTask{AssigneeId: userId, WeekDay: weekday})
	if err != nil {
		return err
	}

	user.SetState(userId, userName, userStates, user.State{Code: user.SCHEDULE_ENTER_TITLE, Request: string(b)})
	return nil
}

func AddScheduleTask(scheduleTask models.ScheduleTask) error {
	bytesRepr, err := json.Marshal(scheduleTask)
	if err != nil {
		return err
	}
	url := DefaultServiceUrl + fmt.Sprintf("%v/schedule/", scheduleTask.AssigneeId)

	_, err = http.Post(url, "application/json", bytes.NewBuffer(bytesRepr))

	if err != nil {
		return err
	}

	return nil
}

func GetWeekdaySchedule(userId int, weekday time.Weekday) ([]models.ScheduleTask, string, error) {
	url := DefaultServiceUrl + fmt.Sprintf("%v/schedule/%s/", userId, weekday.String())

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
	fmt.Fprintf(&output, "%s:\n\n", services.WeekdayToStr(weekday))
	for i, scheduleTask := range weekdaySchedule {
		fmt.Fprintf(&output, "Задача %v\n", i+1)
		fmt.Fprintf(&output, "%s %s\n", emojiTitle, scheduleTask.Title)
		fmt.Fprintf(&output, "%s %s - %s \n", emojiTime, scheduleTask.Start.Format(layoutTime), scheduleTask.End.Format(layoutTime))
		fmt.Fprintf(&output, "%s %s\n", emojiPlace, scheduleTask.Place)
		fmt.Fprintf(&output, "%s %s\n\n", emojiSpeaker, scheduleTask.Speaker)
	}

	return weekdaySchedule, output.String(), nil
}

func ClearAll(userId int) error {
	_, err := utils.Delete(DefaultServiceUrl + fmt.Sprintf("%v/schedule/", userId))
	if err != nil {
		return err
	}
	return nil
}
