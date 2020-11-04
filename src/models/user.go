package models

type User struct {
	Id         int    `json:"user_id"`
	UserName   string `json:"username"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
}
