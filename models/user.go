package models

import "database/sql"

type User struct{
	User_id int 
	Name string
	Email string
}

var Users = map[int]*User{}
var Seq = 1

func Conn()(*sql.DB, error){
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_api")
	if err != nil{
		return nil, err
	}
	return db, nil
	
}
