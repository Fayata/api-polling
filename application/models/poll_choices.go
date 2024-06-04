package models

import (

	"gorm.io/gorm"
)

type PollChoice struct {
	gorm.Model
	PollID int    `gorm:"column:poll_id"`
	Option string `gorm:"column:option"`
}

type UserChoices struct {
	gorm.Model
	ChoiceID uint `gorm:"column:choice_id"`
	UserID   uint `gorm:"column:user_id"`
	PollID   int  `gorm:"column:poll_id"`
}

type PollingResult struct {
	ChoiceID  uint `json:"choice_id"`
	VoteCount int  `json:"vote_count"`
}


