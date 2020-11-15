package models

import "time"

type ScheduleTask struct {
	Id         int          `json:"id"`
	AssigneeId int          `json:"assignee_id"`
	WeekDay    time.Weekday `json:"week_day"`
	Title      string       `json:"title"`
	Place      string       `json:"place"`
	Speaker    string       `json:"speaker"`
	Start      time.Time    `json:"begin"`
	End        time.Time    `json:"end"`
}
