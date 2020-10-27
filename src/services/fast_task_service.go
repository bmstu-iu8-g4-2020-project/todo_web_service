package services

import (
	"todo_web_service/src/models"
)

type DatastoreFastTask interface {
	AddFastTask(fastTask models.FastTask) error
	GetFastTask(fastTask models.FastTask) (models.FastTask, error)
	UpdateFastTask(fastTask models.FastTask) error
	DeleteFastTask(fastTask models.FastTask) error
}

// func (db *DataBase) AddFastTask(fastTask models.FastTask) error

// func (db *DataBase) GetFastTask(fastTask models.FastTask) (models.FastTask, error)

// func (db *DataBase) UpdateFastTask(fastTask models.FastTask) error

// func (db *DataBase) DeleteFastTask(fastTask models.FastTask) error
