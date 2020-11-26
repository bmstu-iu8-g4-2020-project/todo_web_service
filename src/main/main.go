package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/handlers"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/services"
)

const (
	pathToScheme = "../database/init_db.sql"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	db, err := services.NewDB(dbName, dbUser, dbPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	services.Setup(pathToScheme, db)
	fmt.Println("Database is ready!")

	env := &handlers.Environment{Db: db}

	r := mux.NewRouter()

	// User.
	r.HandleFunc("/user/", env.AddUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/user/", env.GetUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", env.GetUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", env.UpdateUserStateHandler).Methods(http.MethodPut)
	// Fast Task.
	r.HandleFunc("/{id}/fast_task/", env.AddFastTaskHandler).Methods(http.MethodPost)
	r.HandleFunc("/fast_task/", env.GetAllFastTasksHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}/fast_task/", env.GetFastTasksHandler).Methods(http.MethodGet)
	r.HandleFunc("/fast_task/", env.UpdateFastTasksHandler).Methods(http.MethodPut)
	r.HandleFunc("/{id}/fast_task/{ft_id}", env.DeleteFastTaskHandler).Methods(http.MethodDelete)
	// Schedule.
	r.HandleFunc("/{id}/schedule/", env.AddScheduleTaskHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}/schedule/{week_day}/", env.GetScheduleTaskHandler).Methods(http.MethodGet)
	r.HandleFunc("/{id}/schedule/", env.UpdateScheduleTaskHandler).Methods(http.MethodPut)
	r.HandleFunc("/{id}/schedule/{sch_id}/", env.DeleteScheduleTaskHandler).Methods(http.MethodDelete)
	r.HandleFunc("/{id}/schedule/delete/{week_day}/", env.DeleteScheduleWeekHandler).Methods(http.MethodDelete)
	r.HandleFunc("/{id}/schedule/", env.ClearAllHandler).Methods(http.MethodDelete)
	// Suburban.
	r.HandleFunc("/suburban", handlers.SuburbanHandler).Methods(http.MethodGet)

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
