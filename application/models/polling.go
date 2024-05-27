package models

import (
	// "api-polling/application/models"
	"api-polling/system/database"
	"log"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Polling struct {
	gorm.Model
	Title    string `gorm:"column:title"`
	Question string `gorm:"column:question"`
	PollID   int    `gorm:"column:poll_id"`
	Type     string `gorm:"column:type"`
	URL      string `gorm:"column:url"`
	Option   string `gorm:"column:option"`
	ImageURL string `gorm:"column:image_url"`
	// Banner   PollingBanner //`gorm:"foreign_key:banner_id"`
	// Choices  []PollChoice  `gorm:"foreign_key:poll_id"`
}

// type PollingBanner struct {
// 	// gorm.Model
// 	PollID int    `gorm:"column:poll_id"`
// 	Type   string `gorm:"column:type"`
// 	URL    string `gorm:"column:url"`
// }

// type PollChoice struct {
// 	// gorm.Model
// 	Option   string `gorm:"column:option"`
// 	PollID   int    `gorm:"foreign_key:poll_id"`
// 	ImageURL string `gorm:"column:image_url"`
// }

///////////////////CMS////////////////////

func (p *Polling) Create() error {
	db := database.GetDB()
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
	existingPolling.Option = p.Option

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&existingPolling).Error; err != nil {
		log.Println("Gagal mengupdate data polling:", err)
		return err
	}

	return nil
}

func (p *Polling) Delete(id int) error {
	db := database.GetDB()
	if err := db.Where("poll_id = ?", id).Delete(&Polling{}).Error; err != nil {
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
	return db.Preload("Choices").Preload("Banner").First(p, id).Error
}

func (up *Polling) GetAll() ([]Polling, error) {
	var polls []Polling
	db := database.GetDB()
	err := db.Find(&polls).Error
	return polls, err
}

func (p *Polling) GetTotalPolls() (int64, error) {
	db := database.GetDB()
	var total int64
	if err := db.Model(&Polling{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// Fungsi untuk memeriksa apakah polling sudah disubmit dan ended
func (p *Polling) CheckPollStatus(e echo.Context) (bool, bool, error) {
	userID := e.Get("user_id").(int)

	var userChoice UserChoice
	err := database.GetDB().Where("user_id = ? AND poll_id = ?", userID, p.ID).First(&userChoice).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, false, err
	}

	isSubmitted := err == nil
	isEnded := false

	return isSubmitted, isEnded, nil
}
