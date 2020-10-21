package samples

import "time"

// Однозначно переводит день недели(с учётом чс/зн) в массив с занятиями на этот день.
type Schedule map[Date][]Cell

type Date struct {
	WeekDay    time.Weekday `json:"weekday"`
	IsEvenWeek bool         `json:"even_week"`
}

type Cell struct {
	Subject   string    `json:"subject"`
	Speaker   string    `json:"speaker"`
	Place     string    `json:"place"`
	Classroom string    `json:"classroom"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
