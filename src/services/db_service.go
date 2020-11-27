package services

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/models"
)

type DataBase struct {
	*sql.DB
}

type Datastore interface {
	// User.
	AddUser(user models.User) error
	GetUsers() ([]models.User, error)
	GetUser(userId int) (models.User, error)
	UpdateState(user models.User) error

	// FastTask.
	AddFastTask(fastTask models.FastTask) error
	GetAllFastTasks() ([]models.FastTask, error)
	GetFastTasks(assigneeId int) ([]models.FastTask, error)
	UpdateFastTasks(fastTasks []models.FastTask) error
	DeleteFastTask(ftId int) error

	// Schedule.
	AddScheduleTask(scheduleTask models.ScheduleTask) error
	GetSchedule(assigneeId int, weekday time.Weekday) ([]models.ScheduleTask, error)
	UpdateScheduleTask(scheduleTask models.ScheduleTask) error
	DeleteScheduleTask(schId int) error
	DeleteScheduleWeek(assigneeId int, weekday time.Weekday) error
	ClearAll(assigneeId int) error
}

func SetDBConfig() string {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	return fmt.Sprintf("host=todo_postgres port=? user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
}

func NewDB(dbSourceName string) (*DataBase, error) {
	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		return nil, err
	}

	// Проверка соединения с бд.
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database!")

	return &DataBase{db}, nil
}

func Setup(pathToFile string, db *DataBase) {
	file, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println("setup file opening error: ", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("error after opening setup file: ", err)
		return
	}

	bytes := make([]byte, stat.Size())
	_, err = file.Read(bytes)
	if err != nil {
		fmt.Println("error after opening setup file: ", err)
		panic(err)
	}

	command := string(bytes)
	_, err = db.Exec(command)
	if err != nil {
		fmt.Println("command error")
	}
}
