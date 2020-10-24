package main

import (
	"database/sql"
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
	connStr       = "user=postgres password=mypass dbname=todo_web_service sslmode=disable"
)

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

type Message struct {
	UserName string `json:"username"`
	ChatID   int64  `json:"chat_id"`
	Text     string `json:"text"`
}

type ScheduleResponse struct {
	Search   Search
	Segments []Segment
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

type DataBase struct {
	*sql.DB
}

type DatastoreUser interface {
	AddUserToDB(user User) error
	UserInfo(user User) (User, error)
}

type EnvUser struct {
	Db DatastoreUser
}

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
}

func (db *DataBase) AddUserToDB(user User) error {
	result, err := db.Exec("INSERT INTO tg_user (username, user_id) values ($1, $2)", user.UserName, user.Id)
	if err != nil {
		return err
	}

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (db *DataBase) UserInfo(user User) (User, error) {
	row := db.QueryRow("SELECT * FROM tg_user WHERE user_id= $1", user.Id)
	err := row.Scan(&user.UserName, &user.Id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (env *EnvUser) AddUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	err = env.Db.AddUserToDB(user)
	if err != nil {
		log.Fatal(err)
	}
}

func (env *EnvUser) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	user, err = env.Db.UserInfo(user)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(user)
}

func main() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	envUser := &EnvUser{Db: db}

	r := mux.NewRouter()

	r.HandleFunc("/", ExampleHandler).Methods("POST")
	r.HandleFunc("/suburban", SuburbanHandler).Methods("GET")

	r.HandleFunc("/user/info", envUser.GetUserInfo).Methods("POST")
	r.HandleFunc("/user", envUser.AddUser).Methods("POST")

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal(err)
	}
}
