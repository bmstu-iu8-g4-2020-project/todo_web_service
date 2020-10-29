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
	UserInfo(user models.User) (models.User, error)

	// FastTask
	GetAllFastTasks() ([]models.FastTask, error)
	AddFastTask(fastTask models.FastTask) error
	//	GetFastTask(fastTask models.FastTask) (models.FastTask, error)
	//	UpdateFastTask(fastTask models.FastTask) error
	//	DeleteFastTask(fastTask models.FastTask) error
}

func NewDB(dbUser string, dbPassword string) (*DataBase, error) {
	dbSourceName := fmt.Sprintf("user=%s password=%s dbname=todownik sslmode=disable", dbUser, dbPassword)

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
