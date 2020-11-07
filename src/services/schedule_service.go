package services

import (
	"fmt"
	"time"
	"todo_web_service/src/models"
)

func ParseStrToWeekday(strWeekday string) (time.Weekday, error) {
	var weekDays = map[string]time.Weekday{
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
		"Sunday":    time.Sunday,
	}
	if weekday, ok := weekDays[strWeekday]; ok {
		return weekday, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", strWeekday)
}

func ParseWeekdayToStr(weekday time.Weekday) string {
	var weekDays = map[time.Weekday]string{
		time.Monday:    "Понедельник",
		time.Tuesday:   "Вторник",
		time.Wednesday: "Среда",
		time.Thursday:  "Четверг",
		time.Friday:    "Пятница",
		time.Saturday:  "Суббота",
		time.Sunday:    "Воскресенье",
	}

	return weekDays[weekday]
}

func (db *DataBase) InitSchedule(assigneeId int) ([]models.Schedule, error) {
	weekDays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}

	for _, weekDay := range weekDays {
		_, err := db.Exec("INSERT INTO schedule (assignee_id, week_day) values ($1, $2);",
			assigneeId, weekDay.String())
		if err != nil {
			return nil, err
		}
	}

	rows, err := db.Query("SELECT * FROM schedule WHERE assignee_id = $1;", assigneeId)
	if err != nil {
		return nil, err
	}

	var assigneeSchedule []models.Schedule
	var schedule models.Schedule
	var strWeekday string
	for rows.Next() {
		err = rows.Scan(&schedule.Id, &schedule.AssigneeId, &strWeekday)
		if err != nil {
			return nil, err
		}

		schedule.WeekDay, err = ParseStrToWeekday(strWeekday)
		if err != nil {
			return nil, err
		}

		assigneeSchedule = append(assigneeSchedule, schedule)
	}

	return assigneeSchedule, nil
}

func (db *DataBase) FillSchedule(scheduleTasks []models.ScheduleTask) error {
	for _, scheduleTask := range scheduleTasks {
		_, err := db.Exec("INSERT INTO schedule_task (schedule_id, title, place, speaker, start_time, end_time) values ($1, $2, $3, $4, $5, $6);",
			scheduleTask.ScheduleId, scheduleTask.Title, scheduleTask.Place,
			scheduleTask.Speaker, scheduleTask.Start, scheduleTask.End)
		if err != nil {
			return err
		}
	}

	return nil
}
