package models

import (
	"api-polling/system/database"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Quiz struct {
	ID       int              `gorm:"column:id"`
	Title    string           `gorm:"column:title"`
	Question string           `gorm:"column:question"`
	Type     string           `gorm:"column:type"`
	URL      string           `gorm:"column:url"`
	Options  []QuizOption     `gorm:"foreignKey:QuizID;references:ID"`
	UserQ    []UserQuizAnswer `gorm:"foreignKey:QuizID;references:ID"`
}

type QuizOption struct {
	ID        int              `gorm:"column:id"`
	Label     string           `gorm:"column:label"`
	QuizID    uint             `gorm:"column:quiz_id"`
	ImageURL  string           `gorm:"column:image_url"`
	IsCorrect bool             `gorm:"column:is_correct"`
	UserQ     []UserQuizAnswer `gorm:"foreignKey:OptionID;references:ID"`
}

type UserQuizAnswer struct {
	ID       int  `gorm:"column:id"`
	UserID   uint `gorm:"column:user_id"`
	QuizID   uint `gorm:"column:quiz_id"`
	OptionID uint `gorm:"column:option_id"`
}

func (q *Quiz) GetByID(id int) error {
	db := database.GetDB()
	return db.Preload("Options").Find(q, id).Error
}

func (q *Quiz) GetAll() ([]Quiz, error) {
	var quizzes []Quiz
	db := database.GetDB()
	err := db.Find(&quizzes).Error
	return quizzes, err
}

func (q *Quiz) GetTotalQuizzes() (int64, error) {
	db := database.GetDB()
	var total int64
	if err := db.Model(&Quiz{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (q *Quiz) CheckQuizStatus(e echo.Context, userID uint) (bool, bool, error) {
    var userAnswer UserQuizAnswer
	err := database.GetDB().Where("user_id = ? AND quiz_id = ?", userID, q.ID).First(&userAnswer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, false, err
	}

	isSubmitted := err == nil
	isEnded := false

	return isSubmitted, isEnded, nil
}

// Fungsi untuk mendapatkan posisi kuis dan total kuis
func (q *Quiz) GetQuizPosition() (int64, int64, error) {
    db := database.GetDB()

    // Hitung total quiz
    var totalQuizzes int64
    if err := db.Model(&Quiz{}).Count(&totalQuizzes).Error; err != nil {
        return 0, 0, err
    }

    // Cari posisi quiz saat ini berdasarkan ID
    var currentQuizPosition int64
    if err := db.Model(&Quiz{}).Where("id < ?", q.ID).Count(&currentQuizPosition).Error; err != nil {
        return 0, 0, err
    }

    return totalQuizzes, currentQuizPosition + 1, nil // +1 karena posisi dimulai dari 1
}