package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	UserName string `json:"username"`
	ChatID   int64  `json:"chat_id"`
	Text     string `json:"text"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	msg.Text = fmt.Sprintf("Hello %s! I'm first testing service.", msg.UserName)

	json.NewEncoder(w).Encode(msg)
}

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
