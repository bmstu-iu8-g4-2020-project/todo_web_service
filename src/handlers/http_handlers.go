package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"todo_web_service/src/models"
)

//func MakeHTTPHandlers(router mux.Router)  {
//	// User handlers. TODO: change nil to Methods.
//	router.HandleFunc("/user/info", nil)
//}

const (
	layoutDate = "2006-01-02"
	layoutTime = "15:04"

	from          = "s2000001" // Kursky Rail Terminal
	to            = "s9601062" // Chekhov
	lang          = "ru_RU"
	transportType = "suburban"
	connStr       = "user=postgres password=mypass dbname=todo_web_service sslmode=disable"
)

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

	sr := &models.ScheduleResponse{}

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