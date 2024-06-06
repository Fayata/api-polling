package cms

import (
	"api-polling/application/entities/dmo/polling"
	"api-polling/system/database"
	"log"

	"gorm.io/gorm"
)

func (r *Repository) Create() error {
	var pollStruct polling.Poll
	db := database.GetDB()
	// Create polling and its choices in a transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pollStruct).Error; err != nil {
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

func (r *Repository) Update(id int) error {
	db := database.GetDB()
	var existingPolling Poll
	if err := db.First(&existingPolling, id).Error; err != nil {
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

func (r *Repository) Delete(id int) error {
	db := database.GetDB()
	if err := db.Where("id = ?", id).Delete(&Poll_Choices{}).Error; err != nil {
		return err
	}

	if err := db.Delete(&Poll{}, id).Error; err != nil {
		return err
	}
	return nil
}
