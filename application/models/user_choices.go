package models

import (
	"api-polling/system/database"
	"errors"
	"log"

	"gorm.io/gorm"
)

type UserChoice struct {
    gorm.Model
    ChoiceID uint `gorm:"not null"`
    UserID   uint `gorm:"not null"`
    PollID   uint `gorm:"not null"`
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
