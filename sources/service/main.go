package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func Pong1(w http.ResponseWriter, r *http.Request) {
	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	msg.Text = "Hello from service: pong1."

	json.NewEncoder(w).Encode(msg)
}

func Pong2(w http.ResponseWriter, r *http.Request) {
	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	msg.Text = "Hello from service: pong2."

	_ = json.NewEncoder(w).Encode(msg)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/pong1", Pong1).Methods("POST")
	r.HandleFunc("/pong2", Pong2).Methods("POST")

	err := http.ListenAndServe(":8001", r)
	if err != nil {
		log.Fatal(err)
	}
}
