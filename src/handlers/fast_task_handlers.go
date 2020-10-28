package handlers

import (
	"encoding/json"
	"net/http"
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

// func (env *Environment) GetFastTask(w http.ResponseWriter, r *http.Request) {}

func (env *Environment) UpdateFastTask(w http.ResponseWriter, r *http.Request) {

}

func (env *Environment) DeleteFastTask(w http.ResponseWriter, r *http.Request) {}
