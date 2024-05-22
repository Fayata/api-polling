package models

import (
	"api-polling/system/database"
	"log"

	"gorm.io/gorm"
)

type PollChoices struct {
	gorm.Model
	PollID int    `gorm:"column:poll_id"`
	Option string `gorm:"column:option"`
}

type UserChoices struct {
	gorm.Model
	ChoiceID int `gorm:"column:choice_id"`
	UserID   int `gorm:"column:user_id"`
	PollID   int  `gorm:"column:poll_id"`
}

type PollingResult struct {
	ChoiceID  int `json:"choice_id"`
	VoteCount int  `json:"vote_count"`
}
type AddPollRequest struct {
    OptionID       int `json:"option_id"`
    QuestionNumber int `json:"question_number"`
}

// Fungsi hasil polling berdasarkan ID polling
func GetPollingResultsByID(pollID uint) ([]map[string]interface{}, error) {
	db := database.GetDB()

	var results []struct {
		PollChoices
		VoteCount int
	}

	err := db.Table("poll_choices pc").
		Select("pc.*, COUNT(uc.id) as vote_count").
		Joins("LEFT JOIN user_choices uc ON pc.id = uc.choice_id").
		Where("pc.poll_id = ?", pollID).
		Group("pc.id").
		Find(&results).Error

	if err != nil {
		log.Println("Failed to fetch results:", err)
		return nil, err
	}

	// Calculate total votes
	var totalVotes int64
	db.Model(&UserChoices{}).Where("poll_id = ?", pollID).Count(&totalVotes)

	formattedResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		percentage := 0.0
		if totalVotes > 0 {
			percentage = float64(result.VoteCount) / float64(totalVotes) * 100
		}

		formattedResults[i] = map[string]interface{}{
			"poll_choice": map[string]interface{}{
				"id":     result.ID,
				"option": result.Option,
			},
			"percentage": percentage,
		}
	}

	// Get polling title
	var polling Polling
	if err := db.Select("title").First(&polling, pollID).Error; err != nil {
		log.Println("Failed to fetch polling title:", err)
		return nil, err
	}

	for i := range formattedResults {
		formattedResults[i]["poll_choice"].(map[string]interface{})["title"] = polling.Title
	}

	return formattedResults, nil
}
