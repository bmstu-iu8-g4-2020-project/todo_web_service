package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"todo_web_service/src/models"
	"todo_web_service/src/services"
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

func (env *Environment) GetScheduleTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assigneeId, _ := strconv.Atoi(params["id"])
	weekday, err := services.StrToWeekday(params["week_day"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	scheduleTasks, err := env.Db.GetSchedule(assigneeId, weekday)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(scheduleTasks)
}