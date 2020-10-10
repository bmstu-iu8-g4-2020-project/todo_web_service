package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

func Controller(r *http.Request) (string, error) {
	// В передаваемом json по полю tag можно
	//пойти либо по pong1 ручке, либо по pong2.

	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		return "", err
	} else {
		return msg.Tag, nil
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	tag, err := Controller(r)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://localhost:8001/%s", tag)

	msg := Message{}
	err = json.NewDecoder(r.Body).Decode(&msg)

	body, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	client := http.Client{}
	resp, _ := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// resp - ответ от сервиса, приводим его в структурку message
	msg = Message{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	_ = json.NewEncoder(w).Encode(msg)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ping", Ping).Methods("GET")

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
