package models

type User struct {
	Id           int    `json:"user_id"`
	UserName     string `json:"username"`
	StateCode    int    `json:"state_code"`
	StateRequest string `json:"state_request"`
}
