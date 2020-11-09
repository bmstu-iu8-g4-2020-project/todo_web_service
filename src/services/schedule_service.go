package services

import (
	"log"
	"time"
	"todo_web_service/src/models"
	"todo_web_service/src/telegram_bot/user"
)

func StrToWeekday(strWeekday string) time.Weekday {
	var weekDays = map[string]time.Weekday{
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
		"Sunday":    time.Sunday,
	}

	return weekDays[strWeekday]
}

func StateCodeToWeekDay(stateCode int) time.Weekday {
	var weekDays = map[int]time.Weekday{
		user.SCHEDULE_FILL_MON: time.Monday,
		user.SCHEDULE_FILL_TUE: time.Tuesday,
		user.SCHEDULE_FILL_WED: time.Wednesday,
		user.SCHEDULE_FILL_THU: time.Thursday,
		user.SCHEDULE_FILL_FRI: time.Friday,
		user.SCHEDULE_FILL_SAT: time.Saturday,
		user.SCHEDULE_FILL_SUN: time.Sunday,
	}

	return weekDays[stateCode]
}

func WeekdayToStr(weekday time.Weekday) string {
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

func (db *DataBase) AddScheduleTask(scheduleTask models.ScheduleTask) error {
	result, err := db.Exec("INSERT INTO schedule (assignee_id, week_day, title, place,"+
		" speaker, start_time, end_time) values ($1, $2, $3, $4, $5, $6, $7);",
		scheduleTask.AssigneeId, scheduleTask.WeekDay, scheduleTask.Title,
		scheduleTask.Place, scheduleTask.Speaker, scheduleTask.Start, scheduleTask.End)

	if err != nil {
		return err
	}

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}
