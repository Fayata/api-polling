package models

import (
	"api-polling/system/database"

	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Quiz struct {
	ID            int            `gorm:"column:id"`
	Name          string         `gorm:"column:name"`
	TotalQuestion string         `gorm:"column:total_question"`
	IsActive      string         `gorm:"column:is_active"`
	StartDate     time.Time      `gorm:"column:start_date"`
	EndDate       time.Time      `gorm:"column:end_date"`
	// UserA         []UserAnswer   `gorm:"foreignKey:QuizID;references:ID"`
	// Question      []QuizQuestion `gorm:"foreignKey:QuizID;references:ID"`
}

func (m *Quiz) TableName() string {
	return "quiz"
}

type QuizQuestionChoice struct {
	ID          int    `gorm:"column:id"`
	Label       string `gorm:"column:label"`
	QuestionID  uint   `gorm:"column:question_id"`
	ChoiceText  string `gorm:"column:choice_text"`
	ChoiceImage string `gorm:"column:choice_image"`
	IsCorrect   bool   `gorm:"column:is_correct"`
}

type QuizQuestion struct {
	ID            int    `gorm:"column:id"`
	QuizID        int    `gorm:"column:quiz_id"`
	Number        int    `gorm:"column:number"`
	QuestionText  string `gorm:"column:question_text"`
	QuestionImage string `gorm:"column:question_image"`
	QuestionURL   string `gorm:"column:question_url"`
}

type UserAnswer struct {
	UserID     int  `gorm:"column:user_id"`
	QuizID     uint `gorm:"column:quiz_id"`
	ChoiceID   uint `gorm:"column:choice_id"`
	QuestionID uint `gorm:"column:question_id"`
}

func (m *UserAnswer) TableName() string {
	return "user_answers"
}

func (q *Quiz) GetByID(id int) (err error) {
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return db.Preload("Options").Find(q, id).Error
	}
	return db.Preload("Options").Find(q, id).Error
}

func (q *Quiz) GetAll() ([]Quiz, error) {
	var quizzes []Quiz
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return quizzes, err
	}
	err = db.Find(&quizzes).Error
	return quizzes, err
}

func (q *Quiz) GetTotalQuizzes() (int64, error) {
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return 0, err
	}
	var total int64
	if err = db.Model(&Quiz{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (q *Quiz) CheckQuizStatus(e echo.Context, userID uint) (bool, bool, error) {
	var userAnswer UserAnswer
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return false, false, err
	}
	err = db.Where("user_id = ? AND quiz_id = ?", userID, q.ID).First(&userAnswer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, false, err
	}

	isSubmitted := err == nil
	isEnded := false

	return isSubmitted, isEnded, nil
}

// Fungsi untuk mendapatkan posisi kuis dan total kuis
func (q *Quiz) GetQuizPosition() (int, error) {
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return 0, err
	}
	// Cari posisi quiz saat ini berdasarkan ID
	var currentQuizPosition int64
	if err = db.Model(&Quiz{}).Where("id < ?", q.ID).Count(&currentQuizPosition).Error; err != nil {
		return 0, err
	}

	return 0, nil
}
