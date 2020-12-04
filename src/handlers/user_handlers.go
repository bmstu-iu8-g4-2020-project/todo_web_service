// Copyright 2020 aaaaaaaalesha <sks2311211@mail.ru>

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/models"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/services"
)

type Environment struct {
	Db services.Datastore
}

func ValidateId(strId string) (int, error) {
	err := validation.Validate(strId, validation.Required, is.Int)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ValidateUserId(strUserId string) (int, error) {
	userId, err := ValidateId(strUserId)
	if err != nil {
		return 0, err
	}

	err = validation.Validate(strUserId, validation.Length(6, 10))
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (env *Environment) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := ValidateUserId(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := env.Db.GetUser(userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		http.Error(w, http.StatusText(http.StatusNotModified), http.StatusNotModified)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Environment) GetUsersHandler(w http.ResponseWriter, _ *http.Request) {
	users, err := env.Db.GetUsers()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
