package services

import (
	"log"
	"time"

	"todo_web_service/src/models"
)

func (db *DataBase) AddFastTask(fastTask models.FastTask) error {
	result, err := db.Exec("INSERT INTO fast_task (assignee_id, chat_id, task_name, notify_interval, deadline) values ($1, $2, $3, $4, $5);",
		fastTask.AssigneeId, fastTask.ChatId, fastTask.TaskName, MillisecondInterval(fastTask.NotifyInterval), fastTask.Deadline)
	if err != nil {
		return err
	}

	log.Println(result.RowsAffected())
	return nil
}

func (db *DataBase) GetAllFastTasks() ([]models.FastTask, error) {
	rows, err := db.Query("SELECT * FROM fast_task;")
	if err != nil {
		return []models.FastTask{}, err
	}

	var fastTasks []models.FastTask
	for rows.Next() {
		fastTask := models.FastTask{}

		var millisecondInterval int64
		err = rows.Scan(&fastTask.Id, &fastTask.AssigneeId, &fastTask.ChatId,
			&fastTask.TaskName, &millisecondInterval, &fastTask.Deadline)
		if err != nil {
			return []models.FastTask{}, err
		}

		fastTask.NotifyInterval = ToDuration(millisecondInterval)

		fastTasks = append(fastTasks, fastTask)
	}

	return fastTasks, nil
}

func (db *DataBase) GetFastTasks(assigneeId int) ([]models.FastTask, error) {
	rows, err := db.Query("SELECT * FROM fast_task WHERE assignee_id = $1;", assigneeId)
	if err != nil {
		return []models.FastTask{}, err
	}

	var fastTasks []models.FastTask
	for rows.Next() {
		fastTask := models.FastTask{}

		var millisecondInterval int64
		err = rows.Scan(&fastTask.Id, &fastTask.AssigneeId, &fastTask.ChatId,
			&fastTask.TaskName, &millisecondInterval, &fastTask.Deadline)
		if err != nil {
			return []models.FastTask{}, err
		}

		fastTask.NotifyInterval = ToDuration(millisecondInterval)

		fastTasks = append(fastTasks, fastTask)
	}

	return fastTasks, nil
}

func (db *DataBase) UpdateFastTasks(fastTasks []models.FastTask) error {
	for _, currTask := range fastTasks {
		newDeadline := currTask.Deadline.Add(currTask.NotifyInterval)

		_, err := db.Exec("UPDATE fast_task SET deadline = $1 WHERE id = $2;",
			newDeadline, currTask.Id)
		if err != nil {
			return err
		}

	}

	return nil
}

func (db *DataBase) DeleteFastTask(ftId int) error {
	_, err := db.Exec("DELETE FROM fast_task WHERE id = $1", ftId)
	if err != nil {
		return err
	}

	return nil
}

func MillisecondInterval(interval time.Duration) int64 {
	return interval.Milliseconds()
}

func ToDuration(millisecondInterval int64) time.Duration {
	return time.Duration(millisecondInterval * 1000000)
}
