package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/gorilla/mux"
	
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/handlers"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/models"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/fast_task"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/user"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/utils"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func AddFastTaskHandler(w http.ResponseWriter, r *http.Request) {
	fastTask := models.FastTask{}
	err := json.NewDecoder(r.Body).Decode(&fastTask)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func TestUser(T *testing.T) {
	NotTestingUrl := utils.DefaultServiceUrl

	ts := InitServer()
	defer ts.Close()

	utils.DefaultServiceUrl = ts.URL + "/"

	err := user.InitUser(123456789, "qwerty12345")
	if err != nil {
		T.Fatal(err)
	}

	// Check userID validation.
	validStrUserId := "123456789"
	invalidStrUserId := "231"

	_, err = handlers.ValidateUserId(validStrUserId)
	if err != nil {
		T.Fatal(err)
	}

	_, err = handlers.ValidateUserId(invalidStrUserId)
	if err == nil {
		T.Fatal()
	}

	utils.DefaultServiceUrl = NotTestingUrl
}

func TestFastTask(T *testing.T) {
	NotTestingUrl := utils.DefaultServiceUrl

	ts := InitServer()
	defer ts.Close()

	utils.DefaultServiceUrl = ts.URL + "/"

	interval, _ := time.ParseDuration("30m")
	err := fast_task.AddFastTask(95495343, "Выпить таблетки", 246676, interval)
	if err != nil {
		T.Fatal(err)
	}

	utils.DefaultServiceUrl = NotTestingUrl
}

var router = mux.NewRouter()

func InitServer() *httptest.Server {
	router.HandleFunc("/user/", AddUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/{id}/fast_task/", AddFastTaskHandler).Methods(http.MethodPost)

	return httptest.NewServer(router)
}
