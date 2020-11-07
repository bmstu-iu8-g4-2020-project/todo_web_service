package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"todo_web_service/src/models"
)

func (env *Environment) InitScheduleHandler(w http.ResponseWriter, r *http.Request) {
	assigneeId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	assigneeSchedule, err := env.Db.InitSchedule(assigneeId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(assigneeSchedule)
}

func (env *Environment) FillScheduleHandler(w http.ResponseWriter, r *http.Request) {
	var scheduleTasks []models.ScheduleTask
	err := json.NewDecoder(r.Body).Decode(&scheduleTasks)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.FillSchedule(scheduleTasks)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
