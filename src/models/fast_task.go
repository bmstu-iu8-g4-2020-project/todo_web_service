package models

type FastTask struct {
	Id         int    `json:"id"`
	AssigneeId int    `json:"assignee_id"`
	Name       string `json:"name"`
	Interval   int64  `json:"interval"`
}
