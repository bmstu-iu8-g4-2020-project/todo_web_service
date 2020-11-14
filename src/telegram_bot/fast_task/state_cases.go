package fast_task

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"strconv"
	"time"
	"todo_web_service/src/models"
	"todo_web_service/src/telegram_bot/user"
	"todo_web_service/src/telegram_bot/utils"
)

func EnterTitle(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) bool {
	var fastTask models.FastTask
	if update.Message.Text == "" {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Нет текстового сообщения, попробуйте ещё раз."))
		return false
	}
	fastTask.TaskName = update.Message.Text
	b, err := json.Marshal(fastTask)
	if err != nil {
		log.Fatal(err)
	}
	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите, с какой периодичностью вам будут приходить сообщения. (Например: 1h10m40s)"))

	user.SetState(update.Message.From.ID, update.Message.From.UserName, userStates,
		user.State{Code: user.FAST_TASK_ENTER_INTERVAL, Request: string(b)})

	return true
}

func EnterInterval(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) bool {
	var fastTask models.FastTask
	interval, err := time.ParseDuration(update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, введённое вами сообщение не удовлетворяет формату. (пример: 1h40m13s) Попробуйте ещё раз."))
		return false
	}
	currUser, err := user.GetUser(update.Message.From.ID)
	if err != nil {
		log.Fatal(err)
	}
	data := []byte(currUser.StateRequest)

	err = json.Unmarshal(data, &fastTask)
	if err != nil {
		log.Fatal(err)
	}
	fastTask.NotifyInterval = interval

	err = AddFastTask(update.Message.From.ID, fastTask.TaskName, update.Message.Chat.ID, fastTask.NotifyInterval)

	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задача успешно добавлена!"))
	user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)

	return true
}

func EnterDeleteNum(update *tgbotapi.Update, bot **tgbotapi.BotAPI, userStates *map[int]user.State) bool {
	fastTasks, _, err := OutputFastTasks(update.Message.From.ID)

	// Считываем порядковый номер задачи, которую нужно удалить.
	ftNumber, err := strconv.Atoi(update.Message.Text)
	if err != nil {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, вы ввели не число. Введите номер задания, который хотите удалить."))
		return false
	}

	if ftNumber <= 0 || ftNumber > len(fastTasks) {
		(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			"Кажется, такого дела не существует. Введите номер задания, который хотите удалить."))
		return false
	}

	fastTaskDeleteUrl := utils.DefaultServiceUrl +
		fmt.Sprintf("%v/fast_task/%v", update.Message.From.ID, fastTasks[ftNumber-1].Id)

	_, err = utils.Delete(fastTaskDeleteUrl)

	if err != nil {
		log.Fatal(err)
	}

	_, output, err := OutputFastTasks(update.Message.From.ID)
	if err != nil {
		log.Fatal(err)
	}

	(*bot).Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Задача %v успешно удалена!\n", ftNumber)+output))
	user.ResetState(update.Message.From.ID, update.Message.From.UserName, userStates)

	return true
}
