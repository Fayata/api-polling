package models

type User struct{
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token"`
}

// var Users = map[int]*User{}
// var Seq = 1 