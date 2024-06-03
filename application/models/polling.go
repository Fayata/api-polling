package models

import (
	// "api-polling/application/models"
	"api-polling/system/database"
	"log"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Polling struct {
	ID       int          `gorm:"column:id"`
	Title    string       `gorm:"column:title"`
	Question string       `gorm:"column:question"`
	Type     string       `gorm:"column:type"`
	URL      string       `gorm:"column:url"`
	ImageURL string       `gorm:"column:image_url"`
	Choices  []PollChoice `gorm:"foreignKey:PollID;references:ID"`
	UserC    []UserChoice `gorm:"foreignKey:PollID;references:ID"`
}

type PollChoice struct {
	ID       int          `gorm:"column:id"`
	Option   string       `gorm:"column:option"`
	PollID   int          `gorm:"column:poll_id"`
	ImageURL string       `gorm:"column:image_url"`
	UserC    []UserChoice `gorm:"foreignKey:ChoiceID;references:ID"`
}

///////////////////CMS////////////////////

func (p *Polling) Create() error {
	db := database.GetDB()
	// Create polling and its choices in a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
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

func (p *Polling) Update(id int) error {
	db := database.GetDB()
	var existingPolling Polling
	if err := db.First(&existingPolling, id).Error; err != nil {
		return err
	}

	existingPolling.Title = p.Title
	existingPolling.Choices = p.Choices

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&existingPolling).Error; err != nil {
		log.Println("Gagal mengupdate data polling:", err)
		return err
	}

	return nil
}

func (p *Polling) Delete(id int) error {
	db := database.GetDB()
	if err := db.Where("id = ?", id).Delete(&PollChoice{}).Error; err != nil {
		return err
	}

	if err := db.Delete(&Polling{}, id).Error; err != nil {
		return err
	}
	return nil
}

////////////////////USERS///////////////////

func (p *Polling) GetByID(id int) error {
	db := database.GetDB()
	return db.Preload("Choices").Find(p, id).Error
}

func (up *Polling) GetAll() ([]Polling, error) {
	var polls []Polling
	db := database.GetDB()
	err := db.Find(&polls).Error
	return polls, err
}

// Fungsi untuk memeriksa apakah polling sudah disubmit dan ended
func (p *Polling) CheckPollStatus(e echo.Context, userId int) (bool, bool, error) {
    var userChoice UserChoice
    err := database.GetDB().Where("user_id = ? AND poll_id = ?", userId, p.ID).Find(&userChoice).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, false, err
	}

	isSubmitted := err == nil
	isEnded := false

	return isSubmitted, isEnded, nil
}
func (pc *PollChoice) GetVoteCount() (float64, error) {
    db := database.GetDB()

    // Hitung total vote untuk poll_id yang sesuai
    var totalVotes int64
    if err := db.Model(&UserChoice{}).Where("poll_id = ?", pc.PollID).Count(&totalVotes).Error; err != nil {
        return 0, err
    }

    // Hitung vote untuk choice_id yang sesuai
    var choiceVotes int64
    if err := db.Model(&UserChoice{}).Where("choice_id = ?", pc.ID).Count(&choiceVotes).Error; err != nil {
        return 0, err
    }

    // Hitung persentase
    var percentage float64
    if totalVotes > 0 {
        percentage = (float64(choiceVotes) / float64(totalVotes)) * 100
    }

    return percentage, nil
}