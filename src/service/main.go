package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	layoutDate = "2006-01-02"
	layoutTime = "15:04"

	from          = "s2000001" // Kursky Rail Terminal
	to            = "s9601062" // Chekhov
	lang          = "ru_RU"
	transportType = "suburban"
)

type Message struct {
	UserName string `json:"username"`
	ChatID   int64  `json:"chat_id"`
	Text     string `json:"text"`
}

type ScheduleResponse struct {
	Search   Search
	Segments []Segment
}

type Search struct {
	To   StationName
	From StationName
}

type StationName struct {
	Title string
}

type Segment struct {
	Arrival   time.Time
	Departure time.Time
	Thread    Thread
}
type Thread struct {
	Number string
	Title  string
}

func SuburbanHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")

	url := fmt.Sprintf("https://api.rasp.yandex.net/v3.0/search/?apikey=%s&format=json&from=%s&to=%s&lang=%s&date=%s&transport_types=%s",
		apiKey, from, to, lang, time.Now().Format(layoutDate), transportType)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	strResp := string(respBody)

	sr := &ScheduleResponse{}

	dataResp := []byte(strResp)

	err = json.Unmarshal(dataResp, sr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, "Три ближайшие электрички в направлении:\n", sr.Search.From.Title, "-->", sr.Search.To.Title)

	counter := 0
	for i := range sr.Segments {
		currSuburban := sr.Segments[i]
		if currSuburban.Departure.After(time.Now()) {
			fmt.Fprintln(w, "Отправляется:", currSuburban.Departure.Format(layoutTime),
				"Прибывает:", currSuburban.Arrival.Format(layoutTime),
				"Поезд:", currSuburban.Thread.Number, currSuburban.Thread.Title)
			counter++
		}

		if counter == 3 {
			break
		}
	}

	if counter == 0 {
		fmt.Fprintln(w, "Ближайших сегодня уже нет :(")
	}

	fmt.Fprintln(w, "\n\nДанные предоставлены сервисом Яндекс.Расписания: http://rasp.yandex.ru/")
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}

	msg.Text = fmt.Sprintf("Hello %s! I'm first testing service.", msg.UserName)

	json.NewEncoder(w).Encode(msg)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", ExampleHandler).Methods("POST")
	r.HandleFunc("/suburban", SuburbanHandler).Methods("GET")

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
