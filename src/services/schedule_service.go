package services

import (
	"errors"
	"time"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/models"
	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot/user"
)

func StrToWeekday(strWeekday string) (time.Weekday, error) {
	switch strWeekday {
	case "Monday", "Понедельник":
		return time.Monday, nil
	case "Tuesday", "Вторник":
		return time.Tuesday, nil
	case "Wednesday", "Среда":
		return time.Wednesday, nil
	case "Thursday", "Четверг":
		return time.Thursday, nil
	case "Friday", "Пятница":
		return time.Friday, nil
	case "Saturday", "Суббота":
		return time.Saturday, nil
	case "Sunday", "Воскресенье":
		return time.Sunday, nil
	}

	return 0, errors.New("the passed string is not a day of the week")
}

func StateCodeToWeekDay(stateCode int) time.Weekday {
	switch stateCode {
	case user.SCHEDULE_FILL_MON:
		return time.Monday
	case user.SCHEDULE_FILL_TUE:
		return time.Tuesday
	case user.SCHEDULE_FILL_WED:
		return time.Wednesday
	case user.SCHEDULE_FILL_THU:
		return time.Thursday
	case user.SCHEDULE_FILL_FRI:
		return time.Friday
	case user.SCHEDULE_FILL_SAT:
		return time.Saturday
	case user.SCHEDULE_FILL_SUN:
		return time.Sunday
	default:
		return 0
	}
}

func WeekdayToStr(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "Понедельник"
	case time.Tuesday:
		return "Вторник"
	case time.Wednesday:
		return "Среда"
	case time.Thursday:
		return "Четверг"
	case time.Friday:
		return "Пятница"
	case time.Saturday:
		return "Суббота"
	case time.Sunday:
		return "Воскресенье"
	default:
		return ""
	}
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
		return nil, err
	}

	var scheduleTasks []models.ScheduleTask
	for rows.Next() {
		scheduleTask := models.ScheduleTask{}
		var strWeekday string

		err = rows.Scan(&scheduleTask.Id, &scheduleTask.AssigneeId, &strWeekday,
			&scheduleTask.Title, &scheduleTask.Place, &scheduleTask.Speaker, &scheduleTask.Start, &scheduleTask.End)

		if err != nil {
			return nil, err
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

func (db *DataBase) UpdateScheduleTask(scheduleTask models.ScheduleTask) error {
	_, err := db.Exec("UPDATE schedule SET week_day = $1, title = $2, place = $3, speaker = $4, start_time = $5, end_time = $6  WHERE id = $7;",
		scheduleTask.WeekDay.String(), scheduleTask.Title, scheduleTask.Place, scheduleTask.Speaker,
		scheduleTask.Start, scheduleTask.End, scheduleTask.Id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) ClearAll(assigneeId int) error {
	_, err := db.Exec("DELETE FROM schedule WHERE assignee_id = $1", assigneeId)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) DeleteScheduleTask(schId int) error {
	_, err := db.Exec("DELETE FROM schedule WHERE id = $1", schId)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) DeleteScheduleWeek(assigneeId int, weekday time.Weekday) error {
	_, err := db.Exec("DELETE FROM schedule WHERE assignee_id = $1 AND week_day = $2",
		assigneeId, weekday.String())
	if err != nil {
		return err
	}

	return nil
}
