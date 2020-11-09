package handlers

import (
	"encoding/json"
	"net/http"
	"todo_web_service/src/models"
)

func (env *Environment) AddScheduleTaskHandler(w http.ResponseWriter, r *http.Request) {
	var scheduleTask models.ScheduleTask
	err := json.NewDecoder(r.Body).Decode(&scheduleTask)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.AddScheduleTask(scheduleTask)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
