package user

import (
	"encoding/json"
	"net/http"

	tgbotapi "github.com/Syfaro/telegram-bot-api"

	"todo_web_service/src/models"
	"todo_web_service/src/telegram_bot/utils"
)

const (
	START = iota // START = 0
	// FastTask.
	FAST_TASK_ENTER_TITLE
	FAST_TASK_ENTER_INTERVAL
	FAST_TASK_DELETE_ENTER_NUM
	// Schedule.
	SCHEDULE_FILL_MON
	SCHEDULE_FILL_TUE
	SCHEDULE_FILL_WED
	SCHEDULE_FILL_THU
	SCHEDULE_FILL_FRI
	SCHEDULE_FILL_SAT
	SCHEDULE_FILL_SUN
	SCHEDULE_ENTER_TITLE
	SCHEDULE_ENTER_PLACE
	SCHEDULE_ENTER_SPEAKER
	SCHEDULE_ENTER_START
	SCHEDULE_ENTER_END
	SCHEDULE_ENTER_OUTPUT_WEEKDAY
	SCHEDULE_UPDATE_ENTER_WEEKDAY
	SCHEDULE_UPDATE_ENTER_NUM_TASK
	SCHEDULE_UPDATE_ENTER_TITLE
	SCHEDULE_UPDATE_ENTER_PLACE
	SCHEDULE_UPDATE_ENTER_SPEAKER
	SCHEDULE_UPDATE_ENTER_START
	SCHEDULE_UPDATE_ENTER_END
	SCHEDULE_DELETE_WEEKDAY_TASK
	SCHEDULE_DELETE_NUM_TASK
	SCHEDULE_DELETE_WEEKDAY
	SCHEDULE_DELETE_CLEARALL
	/* Weather */
	WEATHER_SEND_LOCATION
)

type State struct {
	Code    int    `json:"code"`
	Request string `json:"request"`
}

type StateFunc func(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]State)

func IsStartState(stateCode int) bool {
	return stateCode == START
}

func SendEnteringNotFinished(bot **tgbotapi.BotAPI, chatId int64) {
	(*bot).Send(tgbotapi.NewMessage(chatId, "Вы не закончили ввод данных. \n"+
		"Если хотите прервать ввод, используйте /reset."))
}

func GetStates(userStates *map[int]State) error {
	resp, err := http.Get(utils.DefaultServiceUrl + "user/")
	if err != nil {
		return err
	}

	var users []models.User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return err
	}

	for _, user := range users {
		(*userStates)[user.Id] = State{user.StateCode, user.StateRequest}
	}

	return nil
}

func SetState(userId int, userName string, userStates *map[int]State, state State) error {
	err := UpdateUser(userId, userName, state.Code, state.Request)
	if err != nil {
		return err
	}

	(*userStates)[userId] = State{Code: state.Code, Request: state.Request}

	return nil
}

func ResetState(userId int, userName string, userStates *map[int]State) error {
	err := UpdateUser(userId, userName, START, "{}")
	if err != nil {
		return err
	}
	(*userStates)[userId] = State{Code: START, Request: "{}"}

	return nil
}
