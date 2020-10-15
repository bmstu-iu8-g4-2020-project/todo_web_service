package samples

import "time"

type Task struct {
	Id               int       `json:"id"`
	Title            string    `json:"title"`             // Название таска
	Description      string    `json:"description"`       // Описание таска.
	NotificationFreq int64     `json:"notification_freq"` // Время, которое спят уведомления
	IsAccomplished   bool      `json:"accomplished"`      // Показывает, выполнена ли таска.
	CreationTime     time.Time `json:"creation_time"`     // Время создания таска.
}
