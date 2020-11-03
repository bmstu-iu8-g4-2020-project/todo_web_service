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

	r.HandleFunc("/user", env.AddUser).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", env.GetUserInfo).Methods(http.MethodGet)

	r.HandleFunc("/{id}/fast_task/", env.AddFastTask).Methods(http.MethodPost)
	r.HandleFunc("/fast_task/", env.GetAllFastTasks).Methods(http.MethodGet)
	r.HandleFunc("/{id}/fast_task/", env.GetFastTasks).Methods(http.MethodGet)
	r.HandleFunc("/fast_task/update", env.UpdateFastTasks).Methods(http.MethodPut)
	r.HandleFunc("/{id}/fast_task/{ft_id}", env.DeleteFastTask).Methods(http.MethodDelete)

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
