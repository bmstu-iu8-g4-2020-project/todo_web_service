package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"todo_web_service/src/handlers"
	"todo_web_service/src/services"
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

	env := &handlers.Environment{Db: db}

	r := mux.NewRouter()

	r.HandleFunc("/suburban", handlers.SuburbanHandler).Methods(http.MethodGet)
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
	r.HandleFunc("/{id}/schedule/init", env.InitScheduleHandler).Methods(http.MethodPost)
	r.HandleFunc("/{id}/schedule/fill", env.FillScheduleHandler).Methods(http.MethodPost)

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
