package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"todo_web_service/src/models"
	"todo_web_service/src/services"
)

type UserEnvironment struct {
	Db services.DatastoreUser
}

func (env *UserEnvironment) AddUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.Db.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (env *UserEnvironment) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	paramsFromURL := mux.Vars(r)
	userId, err := strconv.Atoi(paramsFromURL["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{
		Id: userId,
	}

	user, err = env.Db.UserInfo(user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
