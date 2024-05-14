package models

import (
	"api-polling/system/database"
	"errors"
	"log"

	"github.com/labstack/echo"
)

type UserChoice struct {
	ID        int `json:"id"`
	Choice_ID int `json:"choice_id"`
	User_id   int `json:"user_id"`
	// Poll_id   int `json:"poll_id"`
}

func (uc *UserChoice) AddPoll(e echo.Context) error {
	db, err := database.Conn()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}
	defer db.Close()

	// Retrieve user_id safely
	userID, ok := e.Get("user_id").(int)
	if !ok {
		return errors.New("user_id not found or invalid")
	}
	uc.User_id = userID

	// Check if the user has already polled
	checkQuery := `SELECT COUNT(*) FROM user_choice WHERE user_id = ?`
	var count int
	err = db.QueryRow(checkQuery, uc.User_id).Scan(&count)
	if err != nil {
		log.Println("Failed to check polling status:", err)
		return err
	}
	if count > 0 {
		return errors.New("user has already polled")
	}

	// Insert poll record
	insertQuery := `INSERT INTO user_choice (choice_id, user_id) VALUES (?, ?)`
	_, err = db.Exec(insertQuery, uc.Choice_ID, uc.User_id)
	if err != nil {
		log.Println("Failed to add poll:", err)
		return err
	}

	return nil // Success
}
