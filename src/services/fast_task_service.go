package services

import (
	"log"
	"todo_web_service/src/models"
)

type DatastoreFastTask interface {
	AddFastTask(fastTask models.FastTask) error
	GetFastTask(fastTask models.FastTask) (models.FastTask, error)
	UpdateFastTask(fastTask models.FastTask) error
	DeleteFastTask(fastTask models.FastTask) error
}

func (db *DataBase) AddFastTask(fastTask models.FastTask) error {
	result, err := db.Exec("INSERT INTO fast_task (task_name, notify_interval) values ($1, $2)", fastTask.TaskName, fastTask.Interval)
	if err != nil {
		return err
	}

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

// func (db *DataBase) GetFastTask(fastTask models.FastTask) (models.FastTask, error)

// func (db *DataBase) UpdateFastTask(fastTask models.FastTask) error

// func (db *DataBase) DeleteFastTask(fastTask models.FastTask) error
