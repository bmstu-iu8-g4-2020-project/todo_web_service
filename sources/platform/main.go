package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Message struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

func Controller(r *http.Request) (string, error) {
	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		return "", err
	} else {
		return msg.Tag, nil
	}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	// делаем запрос на сервис (по идее надо это как то все более красиво обернуть но суть в том что платформа должна передать
	// запрос от клиента нужному сервису в параметрах урл сервиса, тип тела сообщения, и собсна тело запроса которое мы получили передаем
	tag, err := Controller(r)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("http://localhost:8001/%s", tag)

	resp, err := http.Post(url, "application/json", r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// resp - ответ от сервиса, приводим его в структурку message
	msg := Message{}
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}
	// отправляем клиенту
	_ = json.NewEncoder(w).Encode(msg)
}

func main() {
	// создаем новый роутер))
	r := mux.NewRouter()

	// описываем хэндлер (грубо говоря если клиент сделает запрос по этому урлу какую функцию надо вызвать)
	r.HandleFunc("/ping", Ping).Methods("GET")

	// запускаем сервер
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
