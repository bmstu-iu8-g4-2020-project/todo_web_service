package services

import (
	"database/sql"
	"fmt"
	"todo_web_service/src/models"
)

type DataBase struct {
	*sql.DB
}

type Datastore interface {
	// User
	AddUser(user models.User) error
	GetUsers() ([]models.User, error)
	GetUser(userId int) (models.User, error)
	UpdateState(user models.User) error

	// FastTask
	AddFastTask(fastTask models.FastTask) error
	GetAllFastTasks() ([]models.FastTask, error)
	GetFastTasks(assigneeId int) ([]models.FastTask, error)
	UpdateFastTasks(fastTasks []models.FastTask) error
	DeleteFastTask(ftId int) error

	// ScheduleTask
	AddScheduleTask(scheduleTask models.ScheduleTask) error
}

func NewDB(dbName string, dbUser string, dbPassword string) (*DataBase, error) {
	dbSourceName := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		return nil, err
	}

	// Check connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database!")
	return &DataBase{db}, nil
}
