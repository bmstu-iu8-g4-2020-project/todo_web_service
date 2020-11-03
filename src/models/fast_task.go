package models

import "time"

type FastTask struct {
	Id             int           `json:"id"`
	AssigneeId     int           `json:"assignee_id"`
	ChatId         int64         `json:"chat_id"`
	TaskName       string        `json:"task_name"`
	NotifyInterval time.Duration `json:"interval"`
	Deadline       time.Time     `json:"deadline"`
}
