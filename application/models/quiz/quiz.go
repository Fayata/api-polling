package models

import (
	"api-polling/system/database"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Quiz struct {
    gorm.Model
    Title    string        `gorm:"column:title"`
    Question string        `gorm:"column:question"`
    Banner   QuizBanner    `gorm:"foreignKey:QuizID;references:ID"`
    Options  []QuizOption  `gorm:"foreignKey:QuizID;references:ID"`
}

type QuizBanner struct {
    gorm.Model
    QuizID uint   `gorm:"column:quiz_id"`
    Type   string `gorm:"column:type"`
    URL    string `gorm:"column:url"`
}

type QuizOption struct {
    gorm.Model
    Label     string `gorm:"column:label"`
    QuizID    uint   `gorm:"column:quiz_id"`
    ImageURL  string `gorm:"column:image_url"`
    IsCorrect bool   `gorm:"column:is_correct"`
}

type UserQuizAnswer struct {
    gorm.Model
    UserID   uint `gorm:"column:user_id"`
    QuizID   uint `gorm:"column:quiz_id"`
    OptionID uint `gorm:"column:option_id"`
}
type AddQuizAnswerRequest struct {
    OptionID       uint `json:"option_id"`
    QuestionNumber uint `json:"question_number"`
}

func (q *Quiz) GetByID(id int) error {
    db := database.GetDB()
    return db.Preload("Options").Preload("Banner").First(q, id).Error
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

func (q *Quiz) CheckQuizStatus(e echo.Context) (bool, bool, error) {
    userID := e.Get("user_id").(int)

    var userAnswer UserQuizAnswer
    err := database.GetDB().Where("user_id = ? AND quiz_id = ?", userID, q.ID).First(&userAnswer).Error
    if err != nil && err != gorm.ErrRecordNotFound {
        return false, false, err
    }

    isSubmitted := err == nil
    isEnded := false

    return isSubmitted, isEnded, nil
}
