package models

import (
	"api-polling/system/database"
	"errors"
	"log"

	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	ID            int       `gorm:"column:id"`
	Name          string    `gorm:"column:name"`
	TotalQuestion string    `gorm:"column:total_question"`
	IsActive      string    `gorm:"column:is_active"`
	StartDate     time.Time `gorm:"column:start_date"`
	EndDate       time.Time `gorm:"column:end_date"`
}

func (m *Quiz) TableName() string {
	return "quiz"
}

type QuizQuestionChoice struct {
	ID           int    `gorm:"column:id" json:"id"`
	QuestionID   uint   `gorm:"column:question_id" json:"question_id"`
	ChoiceText   string `gorm:"column:choice_text" json:"label"`
	ChoiceImage  string `gorm:"column:choice_image" json:"image_url"`
	QuestionType string `gorm:"-" json:"-"`
	IsCorrect    bool   `gorm:"column:is_correct" json:"-"`
}

type QuizQuestion struct {
	ID            int    `gorm:"column:id"`
	QuizID        int    `gorm:"column:quiz_id"`
	Number        int    `gorm:"column:number" json:"order"`
	QuestionText  string `gorm:"column:question_text"`
	QuestionImage string `gorm:"column:question_image"`
	// QuestionURL   string `gorm:"column:question_url"`
}

type UserAnswer struct {
	UserID     int  `gorm:"column:user_id"`
	QuizID     int  `gorm:"column:quiz_id"`
	ChoiceID   uint `gorm:"column:choice_id"`
	QuestionID int  `gorm:"column:question_id"`
}

func (m *UserAnswer) TableName() string {
	return "user_answers"
}

// Function for get all questions by quiz id
func GetQuestionByQuizId(id int) (data []QuizQuestion, err error) {
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return data, err
	}
	err = db.Raw("SELECT quiz_id FROM quiz_questions WHERE quiz_id = ?", id).Order("number asc").Scan(&data).Error
	if err != nil {
		log.Println("Failed to fetch question", err)
		return data, err
	}
	return data, err
}

// Function for get all choices by question id
func GetChoiceByQuestionId(ID int) (data []QuizQuestionChoice, err error) {
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return data, err
	}
	err = db.Raw("SELECT question_id FROM quiz_question_choices WHERE question_id = ?", ID).Order("sorting asc").Scan(&data).Error
	if err != nil {
		log.Println("Failed to fetch Option", err)
		return data, err
	}
	return data, err
}

// Function for get all quiz
func (q *Quiz) GetAll() ([]Quiz, error) {
	var quizzes []Quiz
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return quizzes, err
	}
	err = db.Find(&quizzes).Error
	return quizzes, err
}

// Function for get question type quiz by question image
func GetQuestionType(questionImage string) (status bool) {
	if questionImage == "" {
		return false
	}
	return true
}

// Function for get choice type quiz by choice image
func GetChoiceTypeQuiz(choiceImage string) string {
	if choiceImage == "" {
		return "text"
	}
	return "image"
}

// Function for add Quiz
func AddQuiz(user_id int, quiz_id int, question_id int, choice_id uint) error {

	var userAnswer = UserAnswer{
		UserID:     user_id,
		QuizID:     quiz_id,
		QuestionID: question_id,
		ChoiceID:   choice_id,
	}

	db, err := database.InitDB().DbQuiz()

	if err != nil {
		return err
	}
	if err := db.Create(userAnswer).Error; err != nil {
		log.Println("Failed to add poll:", err)
		return err
	}

	return nil
}

func (uc *UserAnswer) AddQuiz() error {
	db, err := database.InitDB().DbPolling()
	if err != nil {
		return err
	}
	// Check if the user has already polled
	var existingVote UserAnswer
	err = db.Raw("SELECT user_id, quiz_id FROM user_answers WHERE user_id = ?, quiz_id = ?", uc.UserID, uc.QuizID).Scan(&existingVote).Error
	if err == nil {
		return errors.New("user has already polled")
	}

	if err := db.Create(uc).Error; err != nil {
		log.Println("Failed to add poll:", err)
		return err
	}

	return nil
}

// Function for check if user submitted the quiz
func IsSubmitted(user_id int, quiz_id int) (status bool, err error) {
	var userAnswer UserAnswer
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return false, err
	}
	err = db.Where("user_id = ? AND quiz_id = ?", user_id, quiz_id).First(&userAnswer).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil
}

// Function for check if quiz is ended(Unfinished)
func IsEnded(ID int) (bool, error) {
	var quiz Quiz
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return false, err
	}
	
	err = db.Where("id = ?", quiz.ID).First(&quiz).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if quiz.EndDate.After(time.Now()) || quiz.StartDate.After(time.Now()) {
		return true, nil
	}
	return false, nil
}

// Function for get total quizzes
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

// Fungsi untuk mendapatkan posisi kuis dan total kuis
func (q *Quiz) GetQuizPosition() (int, error) {
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return 0, err
	}

	// Fetch quiz questions with their choices
	var questions []QuizQuestion
	err = db.Raw("SELECT quiz_id FROM quiz_question WHERE quiz_id = ?", q.ID).Scan(&questions).Error
	if err != nil {
		return 0, err
	}
	// Determine the current question based on whether it's been answered by the user
	currentQuestion := 1
	for _, question := range questions {
		var userAnswer UserAnswer
		err = db.Raw("SELECT user_id, question_id FROM user_answers WHERE user_id = ? AND question_id = ?", q.ID, question.ID).Scan(&userAnswer).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return 0, err
		}
		if err == gorm.ErrRecordNotFound {
			break // This is the current question
		}
		currentQuestion++
	}

	return currentQuestion, nil
}
