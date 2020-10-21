package models

type User struct {
	Email      string `json:"email"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`
}
