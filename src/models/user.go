package models

type UserSignUp struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email      string `json:"email"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`
}
