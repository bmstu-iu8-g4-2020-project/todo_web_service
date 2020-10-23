package models

type Task struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Deadline     int64  `json:"deadline"`
	ReminderTime int64  `json:"reminder_time"`
}
