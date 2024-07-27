package models

import (
	"api-polling/system/database"

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

var dbq *gorm.DB

func init() {
	pollingDB, quizDB := database.InitDB()
	db = pollingDB.Db
	dbq = quizDB.Db
}

// Function for get all questions by quiz id
func GetQuestionByQuizId(id int) (data []QuizQuestion, err error) {
	err = db.Raw("SELECT quiz_id FROM quiz_questions WHERE quiz_id = ?", id).Order("number asc").Scan(&data).Error
	if err != nil {
		log.Println("Failed to fetch question", err)
		return data, err
	}
	return data, err
}

// Function for get all choices by question id
func GetChoiceByQuestionId(ID int) (data []QuizQuestionChoice, err error) {
	err = dbq.Raw("SELECT question_id FROM quiz_question_choices WHERE question_id = ?", ID).Order("sorting asc").Scan(&data).Error
	if err != nil {
		log.Println("Failed to fetch Option", err)
		return data, err
	}
	return data, err
}

// Function for get all quiz
func (q *Quiz) GetAll() (quizzes []Quiz, err error) {

	err = dbq.Find(&quizzes).Error
	return quizzes, err
}

// Function for get question type quiz by question image
func GetQuestionType(questionImage string) (status bool) {

	var count QuizQuestion
	err := dbq.Raw("SELECT * FROM quiz_questions WHERE question_image = ?", questionImage).Scan(&count).Error
	if err != nil {
		return false
	}
	if questionImage == "" {
		return false
	}
	return true
}

// Function for get choice type quiz by choice image
func GetChoiceTypeQuiz(choiceImage string) string {

	var count QuizQuestionChoice
	err := dbq.Raw("SELECT *FROM quiz_question_choices   WHERE choice_image = ?", choiceImage).Scan(&count).Error
	if err != nil {
		return "text"
	}
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

	if err := dbq.Create(userAnswer).Error; err != nil {
		log.Println("Failed to add poll:", err)
		return err
	}

	return nil
}

// Function for check if user submitted the quiz
func IsSubmitted(userID int, quizID int) (bool, error) {
	var ua UserAnswer
	var count int
	err := dbq.Raw("SELECT * FROM user_answers WHERE user_id = ? AND quiz_id = ?", ua.UserID, ua.QuizID).Scan(&count).Error
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

	err := dbq.Raw("SELECT * FROM quiz WHERE id = ?", quiz.ID).Scan(&quiz).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	waktuSaatIni := time.Now()

    if waktuSaatIni.After(quiz.StartDate) && waktuSaatIni.Before(quiz.EndDate) {
        return true, nil
    }
	return false, nil

}

// // Function for get total quizzes
// func (q *Quiz) GetTotalQuizzes() (int64, err error) {

// 	if err != nil {
// 		return 0, err
// 	}
// 	var total int64
// 	if err = db.Model(&Quiz{}).Count(&total).Error; err != nil {
// 		return 0, err
// 	}
// 	return total, nil
// }

// Fungsi untuk mendapatkan posisi kuis dan total kuis
// func (q *Quiz) GetQuizPosition() (int, err error) {

// 	err = db.Raw("SELECT quiz_id FROM quiz_question WHERE quiz_id = ?", q.ID).Scan(&questions).Error
// 	if err != nil {
// 		return 0, err
// 	}
// 	// Determine the current question based on whether it's been answered by the user
// 	currentQuestion := 1
// 	for _, question := range questions {
// 		var userAnswer UserAnswer
// 		err = db.Raw("SELECT user_id, question_id FROM user_answers WHERE user_id = ? AND question_id = ?", q.ID, question.ID).Scan(&userAnswer).Error
// 		if err != nil && err != gorm.ErrRecordNotFound {
// 			return 0, err
// 		}
// 		if err == gorm.ErrRecordNotFound {
// 			break // This is the current question
// 		}
// 		currentQuestion++
// 	}

// 	return currentQuestion, nil
// }
