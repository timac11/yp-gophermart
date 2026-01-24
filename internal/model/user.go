package model

type User struct {
	Id       string
	Login    string
	Password string
}

type UserLoginDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
