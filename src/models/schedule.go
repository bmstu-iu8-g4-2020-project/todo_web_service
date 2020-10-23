package models

type Schedule struct {
	Id         int    `json:"id"`
	AssigneeId int    `json:"assignee_id"`
	DayOfWeek  string `json:"day_of_week"`
}
