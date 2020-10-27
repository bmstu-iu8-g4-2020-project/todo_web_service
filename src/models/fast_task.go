package models

import "time"

type FastTask struct {
	Id         int           `json:"id"`
	AssigneeId int           `json:"assignee_id"`
	TaskName   string        `json:"task_name"`
	Interval   time.Duration `json:"interval"`
}
