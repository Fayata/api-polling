package models

import (
	"api-polling/system/database"
	"log"
)

type User struct {
	ID       int              `gorm:"column:id"`
	Username string           `gorm:"column:username"`
	Email    string           `gorm:"column:email"`
	Password string           `gorm:"column:password"`
	Token    string           `gorm:"column:token"`
	UserC    []UserChoice     `gorm:"foreignKey:UserID;references:ID"`
	UserQ    []UserAnswer `gorm:"foreignKey:UserID;references:ID"`
}

func (m *User) TableName() string {
	return "user"
}


func (u *User) Login(email string, password string) (*User, error) {
	db := database.GetDB()
	var user User
	err := db.Where("email = ? AND password = ?", email, password).First(&user).Error
	if err != nil {
		log.Println("Gagal melakukan login:", err)
		return nil, err
	}
	return &user, nil
}

func (u *User) IsTokenValid(token string) bool {
	db := database.GetDB()
	err := db.Where("id = ? AND token = ?", u.ID, token).First(u).Error
	return err == nil
}

// func (u *User) Register() error {
// 	db, err := database.Conn()
// 	if err != nil {
// 		log.Println("Failed to connect to database:", err)
// 		return err
// 	}
// 	defer db.Close()

// 	// Hash the password
// 	hashedPassword, err := HashPassword(u.Password)
// 	if err != nil {
// 		log.Println("Error hashing password:", err)
// 		return err
// 	}

// 	query := "INSERT INTO user (username, email, password, token) VALUES (?, ?, ?, ?)"
// 	err = db.QueryRow(query, u.Username, u.Email, hashedPassword, u.Token).Scan(&u.ID)
// 	if err != nil {
// 		log.Println("Registration query failed:", err)
// 		return err
// 	}
// 	return nil
// }

// // HashPassword function (using bcrypt)
// func HashPassword(password string) (string, error) {
// 	hash, err := md5.Generate([]byte(password), 10)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(hash), nil
// }
