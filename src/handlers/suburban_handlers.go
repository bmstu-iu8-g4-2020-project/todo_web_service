package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"todo_web_service/src/models"
)

const (
	layoutDate = "2006-01-02"
	layoutTime = "15:04"

	from          = "s9601062" // Kursky Rail Terminal
	to            = "s2000001" // Chekhov
	lang          = "ru_RU"
	transportType = "suburban"
)

func SuburbanHandler(w http.ResponseWriter, _ *http.Request) {
	apiKey := os.Getenv("API_KEY")

	url := fmt.Sprintf("https://api.rasp.yandex.net/v3.0/search/?apikey=%s&format=json&from=%s&to=%s&lang=%s&date=%s&transport_types=%s",
		apiKey, from, to, lang, time.Now().Format(layoutDate), transportType)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	strResp := string(respBody)

	sr := &models.ScheduleResponse{}

	dataResp := []byte(strResp)

	err = json.Unmarshal(dataResp, sr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintln(w, "Три ближайшие электрички в направлении:\n",
		sr.Search.From.Title, "-->", sr.Search.To.Title)

	counter := 0
	for i := range sr.Segments {
		currSuburban := sr.Segments[i]
		if currSuburban.Departure.After(time.Now()) {
			_, _ = fmt.Fprintln(w, "Отправляется:", currSuburban.Departure.Format(layoutTime),
				"Прибывает:", currSuburban.Arrival.Format(layoutTime),
				"Поезд:", currSuburban.Thread.Number, currSuburban.Thread.Title)
			counter++
		}

		if counter == 3 {
			break
		}
	}

	if counter == 0 {
		_, _ = fmt.Fprintln(w, "Ближайших сегодня уже нет :(")
	}

	_, _ = fmt.Fprintln(w, "\n\nДанные предоставлены сервисом Яндекс.Расписания: http://rasp.yandex.ru/")
}
