package models

import (
	"api-polling/system/database"
	"log"

	"gorm.io/gorm"
)

type Polling struct {
	gorm.Model
	Title   string       `gorm:"column:title"`
	Choices []PollChoice `gorm:"foreignKey:poll_id"`
}

type PollChoice struct {
	gorm.Model
	Option string `gorm:"column:option"`
	PollID int `gorm:"column:poll_id"`
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
	existingPolling.Choices = p.Choices // Update associated choices

	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&existingPolling).Error; err != nil {
		log.Println("Gagal mengupdate data polling:", err)
		return err
	}

	return nil
}

func (p *Polling) Delete(id int) error {
	db := database.GetDB()
	if err := db.Where("poll_id = ?", id).Delete(&PollChoice{}).Error; err != nil {
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
	return db.Preload("Choices").First(p, id).Error // Menggunakan Preload untuk memuat pilihan bersamaan
}

func (up *Polling) GetAll() ([]Polling, error) {
	var polls []Polling
	db := database.GetDB()
	err := db.Find(&polls).Error // Mengambil semua polling
	return polls, err
}
