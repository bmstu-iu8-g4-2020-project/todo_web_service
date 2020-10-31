package fast_task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo_web_service/src/models"
)

const (
	DefaultServiceUrl = "http://localhost:8080/"

	emojiAttention = "📢: "
	emojiFastTask  = "⭕ "
)

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
