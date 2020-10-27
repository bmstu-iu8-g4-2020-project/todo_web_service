package models

type User struct {
	Id       int    `json:"user_id"`
	UserName string `json:"username"`
}

type Message struct {
	UserName string `json:"username"`
	ChatID   int64  `json:"chat_id"`
	Text     string `json:"text"`
}
