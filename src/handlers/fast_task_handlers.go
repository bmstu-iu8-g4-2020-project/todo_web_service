package handlers

import (
	"encoding/json"
	"net/http"
	"todo_web_service/src/models"
	"todo_web_service/src/services"
)

type FastTaskEnvironment struct {
	Db services.DatastoreFastTask
}

func (env *FastTaskEnvironment) AddFastTask(w http.ResponseWriter, r *http.Request) {
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

func (env *FastTaskEnvironment) GetFastTask(w http.ResponseWriter, r *http.Request) {}

func (env *FastTaskEnvironment) UpdateFastTask(w http.ResponseWriter, r *http.Request) {}

func (env *FastTaskEnvironment) DeleteFastTask(w http.ResponseWriter, r *http.Request) {}
