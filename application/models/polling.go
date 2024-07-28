package models

import (
	// "api-polling/application/models"
	"api-polling/system/database"
	"errors"

	"log"
	"time"

	"gorm.io/gorm"
)

type Poll struct {
	ID             int    `gorm:"column:id"`
	Title          string `gorm:"column:title"`
	Question_text  string `gorm:"column:question_text"`
	Question_image string `gorm:"column:question_image"`
	ImageURL       string `gorm:"column:image_url"`
	Options_type   string `gorm:"column:options_type"`
	// Banner_type    string     `gorm:"column:banner_type"`
	Start_date time.Time `gorm:"column:start_date"`
	End_date   time.Time `gorm:"column:end_date"`
}

func (m *Poll) TableName() string {
	return "poll"
}

type Poll_Choices struct {
	ID           int    `gorm:"column:id" json:"id"`
	PollID       int    `gorm:"column:poll_id" json:"poll_id"`
	Choice_image string `gorm:"column:choice_image" json:"choice_image"`
	Choice_text  string `gorm:"column:choice_text" json:"choice_text"`
}

type Poll_Result struct {
	Poll_id   int `gorm:"column:poll_id"`
	Choice_id int `gorm:"column:choice_id"`
}

func (m *Poll_Result) TableName() string {
	return "poll_result"
}

type User_Answer struct {
	ID        int `gorm:"column:id"`
	User_Id   int `gorm:"column:user_id"`
	Choice_id int `gorm:"column:choice_id"`
	Poll_Id   int `gorm:"column:poll_id"`
}

func (m *User_Answer) TableName() string {
	return "user_answer"
}

var db *gorm.DB

func init() {
	pollingDB, quizDB := database.InitDB()
	db = pollingDB.Db
	dbq = quizDB.Db
}

///////////////////CMS////////////////////

// func (p *Poll) Create() error {

// 	var
// 	// Create polling and its choices in a transaction
// 	err = db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Create(p).Error; err != nil {
// 			return err
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		log.Println("Gagal membuat polling baru:", err)
// 		return err
// 	}

// 	return nil
// }

// func (p *Poll) Update(id int) error {
// 	db, err := database.GetDB("polling")
// 	if err != nil {
// 		return err
// 	}
// 	var existingPolling Poll
// 	if err = db.First(&existingPolling, id).Error; err != nil {
// 		return err
// 	}

// 	existingPolling.Title = p.Title
// 	// existingPolling.Choices = p.Choices

// 	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&existingPolling).Error; err != nil {
// 		log.Println("Gagal mengupdate data polling:", err)
// 		return err
// 	}

// 	return nil
// }

// func (p *Poll) Delete(id int) error {
// 	db, err := database.GetDB("polling")
// 	if err != nil {
// 		return err
// 	}
// 	if err = db.Where("id = ?", id).Delete(&Poll_Choices{}).Error; err != nil {
// 		return err
// 	}

// 	if err := db.Delete(&Poll{}, id).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

////////////////////USERS///////////////////

func (p *Poll) GetByID(id int) (err error) {

	if err != nil {
		return db.First(p, id).Error
	}
	return db.First(p, id).Error
}

func (up *Poll) GetAll() ([]Poll, error) {
	var polls []Poll

	err := db.Find(&polls).Error // Mengambil semua polling
	return polls, err
}

// Fungsi untuk memeriksa apakah polling sudah disubmit dan ended
func IsSubmittedPoll(User_Id int) (status bool, err error) {
	var userAnswer User_Answer
	err = db.Where("user_id = ?", userAnswer.User_Id).First(&userAnswer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil

}
func IsEndedPoll(ID int) (bool, error) {
	var poll Poll

	err := db.Raw("SELECT id FROM poll WHERE id = ?", poll.ID).Scan(&poll).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	waktuSaatIni := time.Now()

	if waktuSaatIni.After(poll.Start_date) && waktuSaatIni.Before(poll.End_date) {
		return true, nil
	}
	return false, nil
}

// Fungsi untuk mendapatkan persentase vote pada pilihan
func (pc *Poll_Choices) GetVotePercentage(poll_id int) (float32, error) {

	// Hitung total vote untuk poll_id yang sesuai
	var totalVotes int64
	if err := db.Model(&User_Answer{}).Where("poll_id = ?", poll_id).Count(&totalVotes).Error; err != nil {
		return 0, err
	}

	// Hitung vote untuk choice_id yang sesuai
	var choiceVotes int64
	if err := db.Model(&User_Answer{}).Where("choice_id = ? AND poll_id = ?", pc.ID, poll_id).Count(&choiceVotes).Error; err != nil {
		return 0, err
	}

	// Hitung persentase
	var percentage float32
	if totalVotes > 0 {
		percentage = (float32(choiceVotes) / float32(totalVotes)) * 100
	}

	return percentage, nil
}

// Function for banner_type polling
func GetBannerType(questionImages string) string {
	var p Poll
	err := db.Raw("SELECT * FROM poll WHERE question_image = ?", questionImages).Scan(&p).Error
	if err != nil {
		log.Println("Failed to fetch poll:", err)
		return "none"
	}
	if p.Question_image != "" {
		return "question"
	} else {
		return "none"
	}
}

// Function for choice_type polling
func GetChoiceType(optionType string) string {

	var count Poll
	err := db.Raw("SELECT * FROM poll WHERE option_type = ?", optionType).Scan(&count).Error
	if err != nil {
		log.Println("Failed to count poll:", err)
		return "text"
	}
	if optionType == "grid" {
		return "image"
	} else if optionType == "list" {
		return "text"
	} else {
		log.Printf("Unexpected option type: %s", optionType)
		return "text"
	}
}

// Fungsi hasil polling berdasarkan ID polling
func GetPollingResultsByID(poll_id int) (_ []map[string]interface{}, err error) {

	// Ambil semua pilihan untuk pollID tertentu
	var pollChoices []Poll_Choices
	err = db.Raw("SELECT * FROM poll_choices WHERE poll_id = ?", poll_id).Scan(&pollChoices).Error // Select all columns
	if err != nil {
		log.Println("Failed to fetch poll choices:", err)
		return nil, err
	}

	// Ambil jumlah vote untuk setiap pilihan
	choiceVotes := make(map[int]int)
	for _, pc := range pollChoices {
		voteCount, err := pc.GetVotePercentage(int(poll_id))
		if err != nil {
			log.Println("Failed to get vote percentage for choice:", err)
			return nil, err
		}
		choiceVotes[pc.ID] = int(voteCount)
	}

	var polling Poll
	if err := db.First(&polling, poll_id).Error; err != nil {
		log.Println("Failed to fetch polling:", err)
		return nil, err
	}

	formattedResults := make([]map[string]interface{}, len(pollChoices))
	for i, pc := range pollChoices {
		formattedResults[i] = map[string]interface{}{
			"poll_choice": map[string]interface{}{ // Add comma after "title"
				"id":           pc.ID,
				"choice_text":  pc.Choice_text,
				"choice_image": pc.Choice_image,
				"title":        polling.Title,
			},
			"percentage": choiceVotes[pc.ID],
		}
	}

	return formattedResults, nil
}

// Function for adding a new poll answer
func (uc *User_Answer) AddPoll() error {

	// Check if the user has already polled
	var existingVote User_Answer
	err := db.Raw("SELECT user_id FROM user_answer WHERE user_id = ? AND poll_id = ?", uc.User_Id, uc.Poll_Id).Scan(&existingVote).Error
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
