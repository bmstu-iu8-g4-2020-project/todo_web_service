package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todo_web_service/src/models"
	"todo_web_service/src/services"
	"todo_web_service/src/telegram_bot/user"
)

const (
	DefaultServiceUrl = "http://localhost:8080/"
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
