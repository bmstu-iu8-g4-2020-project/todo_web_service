package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) GetAllFastTasksHandler(w http.ResponseWriter, _ *http.Request) {
	fastTasks, err := env.Db.GetAllFastTasks()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(fastTasks)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (env *Environment) GetFastTasksHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	assigneeId, err := ValidateUserId(paramFromURL["id"])

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fastTasks, err := env.Db.GetFastTasks(assigneeId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(fastTasks)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		http.Error(w, http.StatusText(http.StatusNotModified), http.StatusNotModified)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) DeleteFastTaskHandler(w http.ResponseWriter, r *http.Request) {
	paramFromURL := mux.Vars(r)
	ftId, err := ValidateId(paramFromURL["ft_id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.DeleteFastTask(ftId)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
