package models

import "time"

type Schedule struct {
	Id         int          `json:"id"`
	AssigneeId int          `json:"assignee_id"`
	WeekDay    time.Weekday `json:"week_day"`
}

type ScheduleTask struct {
	Id         int       `json:"id"`
	ScheduleId int       `json:"schedule_id"`
	Title      string    `json:"title"`
	Place      string    `json:"place"`
	Speaker    string    `json:"speaker"`
	Start      time.Time `json:"begin"`
	End        time.Time `json:"end"`
}
