package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"todo_web_service/src/models"
)

func (env *Environment) AddFastTask(w http.ResponseWriter, r *http.Request) {
	fastTask := models.FastTask{}
	err := json.NewDecoder(r.Body).Decode(&fastTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.Db.AddFastTask(fastTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (env *Environment) GetAllFastTasks(w http.ResponseWriter, r *http.Request) {
	fastTasks, err := env.Db.GetAllFastTasks()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(fastTasks)
}

func (env *Environment) GetFastTasks(w http.ResponseWriter, r *http.Request) {
	assigneeId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fastTasks, err := env.Db.GetFastTasks(assigneeId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(fastTasks)
}

func (env *Environment) UpdateFastTasks(w http.ResponseWriter, r *http.Request) {
	var fastTasks []models.FastTask
	err := json.NewDecoder(r.Body).Decode(&fastTasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.Db.UpdateFastTasks(fastTasks)

	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}
}

func (env *Environment) DeleteFastTask(w http.ResponseWriter, r *http.Request) {}
