package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todo_web_service/src/models"
)

const (
	DefaultServiceUrl = "http://localhost:8080/"
)

func InitScheduleTable(assigneeId int) ([]models.Schedule, error) {
	url := DefaultServiceUrl + fmt.Sprintf("%v/schedule/init", assigneeId)

	resp, err := http.Post(url, "application/json", nil)

	if err != nil {
		return []models.Schedule{}, err
	}

	var assigneeSchedule []models.Schedule

	json.NewDecoder(resp.Body).Decode(&assigneeSchedule)

	return assigneeSchedule, nil
}

func FillSchedule(assigneeId int, scheduleTasks []models.ScheduleTask) error {
	bytesRepr, err := json.Marshal(scheduleTasks)
	if err != nil {
		return err
	}

	url := DefaultServiceUrl + fmt.Sprintf("%v/schedule/fill", assigneeId)

	_, err = http.Post(url, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		return err
	}

	return nil
}
