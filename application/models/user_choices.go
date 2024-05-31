package models

import (
	"api-polling/system/database"
	"errors"
	"log"
)

type UserChoice struct {
	ID       int  `gorm:"column:id"`
	ChoiceID uint `gorm:"not null;column:choice_id"`
	UserID   int `gorm:"not null;column:user_id"`
	PollID   int  `gorm:"not null;column:poll_id"`
}

func (uc *UserChoice) AddPoll() error {
	db := database.GetDB()

	// Check if the user has already polled
	var existingVote UserChoice
	err := db.Where("user_id = ? AND poll_id = ?", uc.UserID, uc.PollID).First(&existingVote).Error
	if err == nil {
		return errors.New("user has already polled")
	}

	// Create a new vote
	if err := db.Create(uc).Error; err != nil {
		log.Println("Failed to add poll:", err)
		return err
	}

	return nil
}
