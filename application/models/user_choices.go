package models


type UserChoice struct {
	ID       int  `gorm:"column:id"`
	ChoiceID uint `gorm:"not null;column:choice_id"`
	UserID   int `gorm:"not null;column:user_id"`
	PollID   int  `gorm:"not null;column:poll_id"`
}
