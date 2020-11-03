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
	"todo_web_service/src/telegram_bot/client"

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
		for _, currTask := range allFastTasks {
			// Если дедлайн "просрочен", отправляем напоминание пользователю
			// и обновляем время следующего дедлайна.
			if time.Now().After(currTask.Deadline) {
				// Отсылаем напоминание пользователю.
				(*bot).Send(tgbotapi.NewMessage(currTask.ChatId, emojiAttention+currTask.TaskName))
				// Добавляем задачу в батч.
				batch = append(batch, currTask)
			}
		}
		if len(batch) != 0 {
			bytesRepr, err := json.Marshal(batch)
			if err != nil {
				log.Fatal(err)
			}
			url := DefaultServiceUrl + FastTaskPostfix

			_, err = client.Put(url, bytes.NewBuffer(bytesRepr))

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

	var output string

	if len(fastTasks) == 0 {
		output = "Дел не нашлось."
		return []models.FastTask{}, output, nil
	}

	output = "Все существующие дела:\n"
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
