package fast_task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Syfaro/telegram-bot-api"

	"todo_web_service/src/models"
	"todo_web_service/src/telegram_bot/utils"
)

func CheckFastTasks(bot **tgbotapi.BotAPI) {
	for {
		var allFastTasks []models.FastTask
		resp, err := http.Get(utils.DefaultServiceUrl + "fast_task/")
		if err != nil {
			log.Fatal(err)
		}
		err = json.NewDecoder(resp.Body).Decode(&allFastTasks)
		if err != nil {
			log.Fatal(err)
		}

		var batch []models.FastTask // Создаём батч для обновления нескольких дедлайнов.
		for _, currTask := range allFastTasks {
			// Если дедлайн "просрочен", отправляем напоминание пользователю
			// и обновляем время следующего дедлайна.
			if time.Now().After(currTask.Deadline) {
				// Отсылаем напоминание пользователю.
				_, _ = (*bot).Send(tgbotapi.NewMessage(currTask.ChatId,
					fmt.Sprintf("%s Напоминание: \n%s", utils.EmojiAttention, currTask.TaskName)))
				// Добавляем задачу в батч.
				batch = append(batch, currTask)
			}
		}
		if len(batch) != 0 {
			bytesRepr, err := json.Marshal(batch)
			if err != nil {
				log.Fatal(err)
			}

			_, err = utils.Put(utils.DefaultServiceUrl+"fast_task/", bytes.NewBuffer(bytesRepr))

			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(time.Second * 1)
	}
}

func OutputFastTasks(assigneeId int) ([]models.FastTask, string, error) {
	var fastTasks []models.FastTask
	fastTaskUrl := fmt.Sprintf("%s%s/fast_task/", utils.DefaultServiceUrl, strconv.Itoa(assigneeId))
	resp, err := http.Get(fastTaskUrl)
	if err != nil {
		return []models.FastTask{}, "", err
	}

	err = json.NewDecoder(resp.Body).Decode(&fastTasks)
	if err != nil {
		return nil, "", err
	}

	var output string

	if len(fastTasks) == 0 {
		output = "Дел не нашлось. Добавьте их с помощью /add_fast_task"
		return []models.FastTask{}, output, nil
	}

	output = "Все существующие дела:\n"
	for i := range fastTasks {
		output += fmt.Sprintf("%s %v) %s (интервал: %s)\n",
			utils.EmojiFastTask, i+1, fastTasks[i].TaskName, fastTasks[i].NotifyInterval.String())
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

	fastTaskUrl := fmt.Sprintf("%s%s/fast_task/", utils.DefaultServiceUrl, strconv.Itoa(userId))

	_, err = http.Post(fastTaskUrl, "application/json", bytes.NewBuffer(bytesRepr))
	if err != nil {
		return err
	}

	return nil
}
