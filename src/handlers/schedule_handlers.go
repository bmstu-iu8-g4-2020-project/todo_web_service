package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"todo_web_service/src/models"
	"todo_web_service/src/services"
)

func ValidateWeekday(weekday string) (time.Weekday, error) {
	err := validation.Validate(weekday, validation.Required,
		validation.In("Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"))

	if err != nil {
		return 0, err
	}
	weekdayTime, err := services.StrToWeekday(weekday)
	if err != nil {
		return 0, err
	}

	return weekdayTime, nil
}

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

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) GetScheduleTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assigneeId, err := ValidateUserId(params["id"])

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	weekday, err := ValidateWeekday(params["week_day"])
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

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) UpdateScheduleTaskHandler(w http.ResponseWriter, r *http.Request) {
	var scheduleTask models.ScheduleTask
	err := json.NewDecoder(r.Body).Decode(&scheduleTask)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.UpdateScheduleTask(scheduleTask)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) ClearAllHandler(w http.ResponseWriter, r *http.Request) {
	assigneeId, err := ValidateUserId(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.ClearAll(assigneeId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) DeleteScheduleTaskHandler(w http.ResponseWriter, r *http.Request) {
	schId, err := ValidateId(mux.Vars(r)["sch_id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.DeleteScheduleTask(schId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) DeleteScheduleWeekHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assigneeId, err := ValidateUserId(params["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	weekday, err := ValidateWeekday(params["week_day"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.DeleteScheduleWeek(assigneeId, weekday)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
