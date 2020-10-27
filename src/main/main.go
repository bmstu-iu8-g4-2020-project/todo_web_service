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

	r.HandleFunc("/", handlers.ExampleHandler).Methods("POST")
	r.HandleFunc("/suburban", handlers.SuburbanHandler).Methods("GET")

	r.HandleFunc("/user", env.AddUser).Methods("POST")
	r.HandleFunc("/user/{id}", env.GetUserInfo).Methods("GET")

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
