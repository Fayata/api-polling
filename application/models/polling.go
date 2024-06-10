package models

import (
	// "api-polling/application/models"
	"api-polling/system/database"
	"errors"
	"log"
	"time"

	"github.com/labstack/echo"
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
	ID           int    `gorm:"column:id"`
	PollID       int    `gorm:"column:poll_id"`
	Choice_image string `gorm:"column:choice_image"`
	Choice_text  string `gorm:"column:choice_text"`
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

///////////////////CMS////////////////////

func (p *Poll) Create() error {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return err
	}
	// Create polling and its choices in a transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("Gagal membuat polling baru:", err)
		return err
	}

	return nil
}

func (p *Poll) Update(id int) error {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return err
	}
	var existingPolling Poll
	if err = db.First(&existingPolling, id).Error; err != nil {
		return err
	}

	existingPolling.Title = p.Title
	// existingPolling.Choices = p.Choices

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&existingPolling).Error; err != nil {
		log.Println("Gagal mengupdate data polling:", err)
		return err
	}

	return nil
}

func (p *Poll) Delete(id int) error {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return err
	}
	if err = db.Where("id = ?", id).Delete(&Poll_Choices{}).Error; err != nil {
		return err
	}

	if err := db.Delete(&Poll{}, id).Error; err != nil {
		return err
	}
	return nil
}

////////////////////USERS///////////////////

func (p *Poll) GetByID(id int) (err error) {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return db.Preload("Choices").Find(p, id).Error
	}
	return db.Preload("Choices").Find(p, id).Error
}

func (up *Poll) GetAll() ([]Poll, error) {
	var polls []Poll
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return polls, err
	}
	err = db.Find(&polls).Error // Mengambil semua polling
	return polls, err
}

// Fungsi untuk memeriksa apakah polling sudah disubmit dan ended
func (p *Poll) CheckPollStatus(e echo.Context, userId int) (bool, bool, error) {
	var userChoice User_Answer
	DB, err := database.InitDB().DbPolling()
	if err != nil {
		return false, false, err
	}
	err = DB.Where("user_id = ? AND poll_id = ?", userId, p.ID).First(&userChoice).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, false, err
	}

	isSubmitted := err == nil
	isEnded := false

	return isSubmitted, isEnded, nil
}

// Fungsi untuk mendapatkan persentase vote pada pilihan
func (pc *Poll_Choices) GetVotePercentage(poll_id int) (int, error) {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return 0, err
	}
	// Hitung total vote untuk poll_id yang sesuai
	var totalVotes int64
	if err = db.Model(&User_Answer{}).Where("poll_id = ?", poll_id).Count(&totalVotes).Error; err != nil {
		return 0, err
	}

	// Hitung vote untuk choice_id yang sesuai
	var choiceVotes int64
	if err := db.Model(&User_Answer{}).Where("choice_id = ? AND poll_id = ?", pc.ID, poll_id).Count(&choiceVotes).Error; err != nil {
		return 0, err
	}

	// Hitung persentase
	var percentage int
	if totalVotes > 0 {
		percentage = (int(choiceVotes) / int(totalVotes)) * 100
	}

	return percentage, nil
}

func (p *Poll) GetBannerType() string {
	if p.Options_type == "list" {
		if p.Question_image != "" {
			return "question"
		} else {
			return "none"
		}
	} else if p.Options_type == "grid" {
		return "image"
	}
	return "none"
}

// Fungsi hasil polling berdasarkan ID polling
func GetPollingResultsByID(poll_id uint) ([]map[string]interface{}, error) {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return nil, err
	}
	// Ambil semua pilihan untuk pollID tertentu
	var pollChoices []Poll_Choices
	if err = db.Where("poll_id = ?", poll_id).Find(&pollChoices).Error; err != nil {
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

	// Ambil data polling berdasarkan ID
	var polling Poll
	if err := db.First(&polling, poll_id).Error; err != nil {
		log.Println("Failed to fetch polling:", err)
		return nil, err
	}

	formattedResults := make([]map[string]interface{}, len(pollChoices))
	for i, pc := range pollChoices {
		formattedResults[i] = map[string]interface{}{
			"poll_choice": map[string]interface{}{
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

func (uc *User_Answer) AddPoll() error {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return err
	}
	// Check if the user has already polled
	var existingVote User_Answer
	err = db.Where("user_id = ? AND poll_id = ?", uc.User_Id, uc.Poll_Id).First(&existingVote).Error
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
