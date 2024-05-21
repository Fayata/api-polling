package models

import (
	"api-polling/system/database"
	"log"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model 
    Username string `gorm:"size:255"`
    Email    string `gorm:"size:255"`
    Password string `gorm:"size:255"`
    Token    string `gorm:"size:255"`
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
