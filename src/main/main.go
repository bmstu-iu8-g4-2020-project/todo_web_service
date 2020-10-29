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

	env := &handlers.Environment{Db: db}

	r := mux.NewRouter()

	r.HandleFunc("/suburban", handlers.SuburbanHandler).Methods("GET")

	r.HandleFunc("/user", env.AddUser).Methods("POST")
	r.HandleFunc("/user/{id}", env.GetUserInfo).Methods("GET")

	r.HandleFunc("/{id}/fast_task/", env.AddFastTask).Methods("POST")
	r.HandleFunc("/fast_task/", env.GetAllFastTasks).Methods("GET")
	//	r.HandleFunc("/{id}/fast_task/", env.GetFastTask).Methods("GET")
	//	r.HandleFunc("/{id}/fast_task/", env.UpdateFastTask).Methods("PUT")
	//	r.HandleFunc("/{id}/fast_task/", env.DeleteFastTask).Methods("DELETE")

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
