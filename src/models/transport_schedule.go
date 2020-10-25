package models

import "time"

type Search struct {
	To   StationName
	From StationName
}

type StationName struct {
	Title string
}

type Segment struct {
	Arrival   time.Time
	Departure time.Time
	Thread    Thread
}

type Thread struct {
	Number string
	Title  string
}

type ScheduleResponse struct {
	Search   Search
	Segments []Segment
}
