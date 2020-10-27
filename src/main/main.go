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
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	db, err := services.NewDB(dbUser, dbPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	envUser := &handlers.UserEnvironment{Db: db}
	envFastTask := &handlers.FastTaskEnvironment{Db: db}

	r := mux.NewRouter()

	r.HandleFunc("/suburban", handlers.SuburbanHandler).Methods("GET")

	r.HandleFunc("/user", envUser.AddUser).Methods("POST")
	r.HandleFunc("/user/{id}", envUser.GetUserInfo).Methods("GET")

	r.HandleFunc("/{id}/fast_task", envFastTask.AddFastTask).Methods("POST")
	r.HandleFunc("/{id}/fast_task/", envFastTask.GetFastTask).Methods("GET")
	r.HandleFunc("/{id}/fast_task/", envFastTask.UpdateFastTask).Methods("PUT")
	r.HandleFunc("/{id}/fast_task/", envFastTask.DeleteFastTask).Methods("DELETE")

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
