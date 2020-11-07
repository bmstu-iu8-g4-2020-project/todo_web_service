package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"todo_web_service/src/models"
)

func (env *Environment) AddFastTaskHandler(w http.ResponseWriter, r *http.Request) {
	fastTask := models.FastTask{}
	err := json.NewDecoder(r.Body).Decode(&fastTask)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.AddFastTask(fastTask)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (env *Environment) GetAllFastTasksHandler(w http.ResponseWriter, r *http.Request) {
	fastTasks, err := env.Db.GetAllFastTasks()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(fastTasks)
}

func (env *Environment) GetFastTasksHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	assigneeId, err := strconv.Atoi(paramFromURL["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fastTasks, err := env.Db.GetFastTasks(assigneeId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(fastTasks)
}

func (env *Environment) UpdateFastTasksHandler(w http.ResponseWriter, r *http.Request) {
	var fastTasks []models.FastTask
	err := json.NewDecoder(r.Body).Decode(&fastTasks)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.UpdateFastTasks(fastTasks)

	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}
}

func (env *Environment) DeleteFastTaskHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	ftId, err := strconv.Atoi(paramFromURL["ft_id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.DeleteFastTask(ftId)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotModified), http.StatusNotModified)
		return
	}
}
