package models

import (
    "api-polling/system/database"
    "log"
)

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token"`
}

func (u *User) Login(email string, password string) (*User, error) {
    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return nil, err
    }
    defer db.Close()

    query := "SELECT id, email, token FROM user WHERE email = ? AND password = ?"
    row := db.QueryRow(query, email, password)

    var user User

    err = row.Scan(&user.ID, &user.Email, &user.Token)
    if err != nil {
        log.Println("Gagal melakukan query:", err)
        return nil, err
    }

    return &user, nil
}
