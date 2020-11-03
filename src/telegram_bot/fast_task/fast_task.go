package fast_task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"todo_web_service/src/models"

	"github.com/Syfaro/telegram-bot-api"
)

const (
	DefaultServiceUrl = "http://localhost:8080/"

	emojiAttention = "📢: "
	emojiFastTask  = "📌 "

	FastTaskPostfix = "fast_task/"
)

func CheckFastTasks(bot **tgbotapi.BotAPI) {
	for {
		var allFastTasks []models.FastTask
		resp, err := http.Get(DefaultServiceUrl + FastTaskPostfix)
		if err != nil {
			log.Fatal(err)
		}
		json.NewDecoder(resp.Body).Decode(&allFastTasks)

		var batch []models.FastTask // Создаём батч для обновления нескольких дедлайнов.
		for i := range allFastTasks {
			currFastTask := allFastTasks[i]
			// Если дедлайн "просрочен", отправляем напоминание пользователю
			// и обновляем время следующего дедлайна.
			if time.Now().After(currFastTask.Deadline) {
				// Отсылаем напоминание пользователю.
				(*bot).Send(tgbotapi.NewMessage(currFastTask.ChatId, emojiAttention+currFastTask.TaskName))
				// Добавляем задачу в батч.
				batch = append(batch, currFastTask)
			}
		}
		if len(batch) != 0 {
			bytesRepr, err := json.Marshal(batch)
			if err != nil {
				log.Fatal(err)
			}
			url := DefaultServiceUrl + FastTaskPostfix + "update"
			_, err = http.Post(url, "application/json", bytes.NewBuffer(bytesRepr))
			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(time.Second * 10)
	}
}

func OutputFastTasks(assigneeId int) ([]models.FastTask, string, error) {
	var fastTasks []models.FastTask
	fastTaskUrl := DefaultServiceUrl + fmt.Sprintf("/%s/fast_task/", strconv.Itoa(assigneeId))
	resp, err := http.Get(fastTaskUrl)
	if err != nil {
		return []models.FastTask{}, "", err
	}

	json.NewDecoder(resp.Body).Decode(&fastTasks)

	output := "Все существующие дела:\n"
	for i := range fastTasks {
		output += emojiFastTask + fmt.Sprintf("%v) %s \n", i+1, fastTasks[i].TaskName)
	}

	return fastTasks, output, nil
}

func AddFastTask(userId int, taskName string, chatID int64, interval time.Duration) error {
	fastTask := models.FastTask{
		AssigneeId:     userId,
		TaskName:       taskName,
		ChatId:         chatID,
		NotifyInterval: interval,
		Deadline:       time.Now().Add(interval),
	}

	bytesRepr, err := json.Marshal(fastTask)
	if err != nil {
		return err
	}

	fastTaskUrl := DefaultServiceUrl + fmt.Sprintf("%s", strconv.Itoa(userId)) + "/fast_task/"

	_, err = http.Post(fastTaskUrl, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		return err
	}

	return nil
}
