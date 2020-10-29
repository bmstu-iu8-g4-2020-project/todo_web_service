package services

import (
	"log"
	"time"
	"todo_web_service/src/models"
)

func (db *DataBase) AddFastTask(fastTask models.FastTask) error {
	millisecondInterval := fastTask.NotifyInterval.Milliseconds()
	result, err := db.Exec("INSERT INTO fast_task (assignee_id, chat_id, task_name, notify_interval, deadline) values ($1, $2, $3, $4, $5)",
		fastTask.AssigneeId, fastTask.ChatId, fastTask.TaskName, millisecondInterval, fastTask.Deadline)
	if err != nil {
		return err
	}

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (db *DataBase) GetAllFastTasks() ([]models.FastTask, error) {
	rows, err := db.Query("SELECT * FROM fast_task")
	if err != nil {
		return []models.FastTask{}, err
	}

	var ftStorage []models.FastTask
	for rows.Next() {
		fastTask := models.FastTask{}

		var millisecondInterval int64
		err = rows.Scan(&fastTask.Id, &fastTask.AssigneeId, &fastTask.ChatId,
			&fastTask.TaskName, &millisecondInterval, &fastTask.Deadline)
		if err != nil {
			return []models.FastTask{}, err
		}

		fastTask.NotifyInterval = time.Duration(millisecondInterval * 1000000)

		ftStorage = append(ftStorage, fastTask)
	}

	return ftStorage, nil
}

// func (db *DataBase) GetFastTask(fastTask models.FastTask) (models.FastTask, error)

// func (db *DataBase) UpdateFastTask(fastTask models.FastTask) error

// func (db *DataBase) DeleteFastTask(fastTask models.FastTask) error
