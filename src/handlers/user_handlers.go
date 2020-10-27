package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"todo_web_service/src/models"
	"todo_web_service/src/services"
)

type Environment struct {
	Db services.Datastore
}

func (env *Environment) AddUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = env.Db.AddUserToDB(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (env *Environment) GetUserInfo(w http.ResponseWriter, r *http.Request) {
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

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	msg := models.Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	msg.Text = fmt.Sprintf("Hello %s! I'm first testing services.", msg.UserName)

	json.NewEncoder(w).Encode(msg)
}
