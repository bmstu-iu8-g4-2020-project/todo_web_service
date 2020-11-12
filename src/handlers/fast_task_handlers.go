package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"todo_web_service/src/models"
)

func ValidateFastTaskId(ftIdStr string) (int, error) {
	err := validation.Validate(ftIdStr, validation.Required, is.Int, validation.Min(0))
	if err != nil {
		return 0, err
	}
	schId, err := strconv.Atoi(ftIdStr)
	if err != nil {
		return 0, err
	}

	return schId, nil
}

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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(fastTasks)
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
	// TODO: удаление не работает, чот криво валидируется.
	ftId, err := ValidateFastTaskId(paramFromURL["ft_id"])
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
