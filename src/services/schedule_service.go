package services

import (
	"errors"
	"fmt"
	"time"
	"todo_web_service/src/models"
	"todo_web_service/src/telegram_bot/user"
)

func StrToWeekday(strWeekday string) (time.Weekday, error) {
	var weekDays = map[string]time.Weekday{
		"Monday":      time.Monday,
		"Tuesday":     time.Tuesday,
		"Wednesday":   time.Wednesday,
		"Thursday":    time.Thursday,
		"Friday":      time.Friday,
		"Saturday":    time.Saturday,
		"Sunday":      time.Sunday,
		"Понедельник": time.Monday,
		"Вторник":     time.Tuesday,
		"Среда":       time.Wednesday,
		"Четверг":     time.Thursday,
		"Пятница":     time.Friday,
		"Суббота":     time.Saturday,
		"Воскресенье": time.Sunday,
	}
	if weekday, ok := weekDays[strWeekday]; ok {
		return weekday, nil
	}

	return 0, errors.New("the passed string is not a day of the week")
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
	_, err := db.Exec("INSERT INTO schedule (assignee_id, week_day, title, place,"+
		" speaker, start_time, end_time) values ($1, $2, $3, $4, $5, $6, $7);",
		scheduleTask.AssigneeId, scheduleTask.WeekDay.String(), scheduleTask.Title,
		scheduleTask.Place, scheduleTask.Speaker, scheduleTask.Start, scheduleTask.End)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) GetSchedule(assigneeId int, weekday time.Weekday) ([]models.ScheduleTask, error) {
	rows, err := db.Query("SELECT * FROM schedule WHERE assignee_id = $1 AND week_day = $2;",
		assigneeId, weekday.String())
	if err != nil {
		return []models.ScheduleTask{}, err
	}

	var scheduleTasks []models.ScheduleTask
	for rows.Next() {
		scheduleTask := models.ScheduleTask{}
		var strWeekday string

		err = rows.Scan(&scheduleTask.Id, &scheduleTask.AssigneeId, &strWeekday,
			&scheduleTask.Title, &scheduleTask.Place, &scheduleTask.Speaker, &scheduleTask.Start, &scheduleTask.End)

		if err != nil {
			return []models.ScheduleTask{}, err
		}

		tempWeekday, err := StrToWeekday(strWeekday)
		if err != nil {
			return nil, err
		}

		scheduleTask.WeekDay = tempWeekday

		scheduleTasks = append(scheduleTasks, scheduleTask)
	}

	return scheduleTasks, nil
}

func (db *DataBase) ClearAll(assigneeId int) error {
	result, err := db.Exec("DELETE FROM schedule WHERE assignee_id = $1", assigneeId)
	if err != nil {
		return err
	}

	fmt.Println(result.RowsAffected())

	return nil
}

func (db *DataBase) DeleteScheduleTask(schId int) error {
	result, err := db.Exec("DELETE FROM schedule WHERE id = $1", schId)
	if err != nil {
		return err
	}

	fmt.Println(result.RowsAffected())

	return nil
}

func (db *DataBase) DeleteScheduleWeek(assigneeId int, weekday time.Weekday) error {
	result, err := db.Exec("DELETE FROM schedule WHERE assignee_id = $1 AND week_day = $2",
		assigneeId, weekday.String())
	if err != nil {
		return err
	}

	fmt.Println(result.RowsAffected())

	return nil
}
