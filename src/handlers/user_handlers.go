package handlers

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"todo_web_service/src/models"
	"todo_web_service/src/services"
)

type Environment struct {
	Db services.Datastore
}

func ValidateUserId(userIdStr string) (int, error) {
	err := validation.Validate(userIdStr, validation.Required, is.Int, validation.Length(6, 10))
	if err != nil {
		return 0, err
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (env *Environment) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = env.Db.AddUser(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (env *Environment) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := ValidateUserId(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := env.Db.GetUser(userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (env *Environment) UpdateUserStateHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = env.Db.UpdateState(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (env *Environment) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := env.Db.GetUsers()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(users)
}
