package models

import (
	"api-polling/system/database"
	"log"
)
type UserChoice struct {
	ID int    `json:"id"`
	Choice_ID int `json:"choice_id"`
	User_id int `json:"user_id"`
}


type PollingResult struct {
	ChoiceID  int `json:"choice_id"`
	VoteCount int `json:"vote_count"`
}

func (uc *UserChoice) AddPoll() error {
    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return err
    }
    defer db.Close()
    query := `
        INSERT INTO user_choice (choice_id, user_id)
        VALUES (?, ?)
    `
    _, err = db.Exec(query, uc.Choice_ID, uc.User_id)
    if err != nil {
        log.Println("Gagal menambahkan polling:", err)
        return err
    }
    return nil
}
