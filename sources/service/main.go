package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Message struct {
	Tag  string    `json:"tag"`
	Text string `json:"text"`
}

func Pong(w http.ResponseWriter, r *http.Request) {
	// декодим в струтру тело запроса
	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	// добавляем тексту "from service."
	msg.Text = msg.Text + " from pong."

	// отправляем ответ
	json.NewEncoder(w).Encode(msg)
}

func Pong1(w http.ResponseWriter, r *http.Request) {
	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	// добавляем тексту "from pong1."
	msg.Text = msg.Text + " from pong1."

	// отправляем ответ
	_ = json.NewEncoder(w).Encode(msg)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/pong", Pong).Methods("POST")
	r.HandleFunc("/pong1", Pong1).Methods("POST")

	err := http.ListenAndServe(":8001", r)
	if err != nil {
		log.Fatal(err)
	}
}

