package models



type User struct{
	User_id int 
	Name string
	Email string
}

var Users = map[int]*User{}
var Seq = 1 