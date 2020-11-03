package services

import (
	"database/sql"
	"fmt"
	"log"
	"todo_web_service/src/models"
)

type DataBase struct {
	*sql.DB
}

type Datastore interface {
	// User
	AddUser(user models.User) error
	UserInfo(userId int) (models.User, error)

	// FastTask
	AddFastTask(fastTask models.FastTask) error
	GetAllFastTasks() ([]models.FastTask, error)
	GetFastTasks(assigneeId int) ([]models.FastTask, error)
	UpdateFastTasks(fastTasks []models.FastTask) error
	DeleteFastTask(ftId int) error
}

func NewDB(dbUser string, dbPassword string) (*DataBase, error) {
	log.Println(dbUser, dbPassword)
	dbSourceName := fmt.Sprintf("user=%s password=%s dbname=todoapp2 sslmode=disable", dbUser, dbPassword)

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
