package models

import (
    "api-polling/system/database"
    "log")
type UserChoice struct {
	ID int    `json:"id"`
	Choice_ID int `json:"choice_id"`
	User_id int `json:"user_id"`
}

type PollingResult struct {
	ChoiceID  int `json:"choice_id"`
	VoteCount int `json:"vote_count"`
}

//models addpolling
func (p *Polling) AddPoll(user_id int, choiceID int) error {
    db, err := database.Conn()
    if err != nil {
        log.Println("Failed to connect to database:", err)
        return err
    }
    defer db.Close()

    query := "INSERT INTO user_choice ( user_id, choice_id) VALUES ( ?, ?)"
    _, err = db.Exec(query, user_id, choiceID)
    if err != nil {
        log.Println("Failed to execute query:", err)
        return err
    }

    return nil
}