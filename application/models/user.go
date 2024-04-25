package models

type User struct{
	User_id int `json:"user_id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token"`
}

// var Users = map[int]*User{}
// var Seq = 1 