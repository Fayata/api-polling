package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type QuestionData struct {
	ID     int `json:"id"`
	QuizID int `json:"quiz_id"`
	Order  struct {
		Total   int `json:"total"`
		Current int `json:"current"`
	} `json:"order"`
	QuestionText string `json:"question_text"`
	Banner       struct {
		ShowBanner bool   `json:"show_banner"`
		URL        string `json:"url"`
	} `json:"banner"`
	Type    string                      `json:"type"`
	Choices []models.QuizQuestionChoice `json:"choices"`
}
var dbq *gorm.DB

func init() {
	pollingDB, quizDB := database.InitDB()
	dbq = pollingDB.Db
	dbq = quizDB.Db
}


func GetQuizByID(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
    }

    var quiz models.Quiz

    // Get quiz by ID
    
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
    }
    err = dbq.First(&quiz, id).Error
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Quiz not found")
    }

	userID := e.Get("user_id").(int)

    // Fetch quiz questions
    var questions []models.QuizQuestion
    err = dbq.Raw("SELECT * FROM quiz_questions WHERE quiz_id = ? ORDER BY number ASC", id).Scan(&questions).Error
    if err != nil {
        log.Println("Failed to get questions:", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
    }

    var questionResults []QuestionData
    for i, question := range questions {
        var questionResult QuestionData
        questionResult.ID = question.ID
        questionResult.QuestionText = question.QuestionText
        questionResult.QuizID = question.QuizID

        // Fetch choices for the current question
        var choiceDataArr []models.QuizQuestionChoice
        err = dbq.Raw("SELECT * FROM quiz_question_choices WHERE question_id = ? ORDER BY sorting ASC", question.ID).Scan(&choiceDataArr).Error
        if err != nil {
            log.Println("Failed to get question choices:", err)
            return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
        }
        
        // Determine question type and assign choices
        questionResult.Choices = choiceDataArr
        for j, choice := range questionResult.Choices {
            choiceDataArr[j].QuestionType = models.GetChoiceType(choice.ChoiceImage)
        }
        questionResult.Type = models.GetChoiceTypeQuiz(questionResult.Choices[0].ChoiceImage) // Assuming all choices have the same type

        questionResult.Order.Total = len(questions)
        questionResult.Order.Current = i + 1

        // Get banner info
        questionResult.Banner.ShowBanner = models.GetQuestionType(question.QuestionImage)
        questionResult.Banner.URL = question.QuestionImage

        questionResults = append(questionResults, questionResult)
    }
	isSubmitted, err := models.IsSubmitted(userID, id)
	if err != nil {
		log.Println("Error checking submission status:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	isEnded, err := models.IsEnded(id) 
	if err!= nil {
        log.Println("Error checking end status:", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
    }

	//  Get meta data
	metaData := database.Meta()

	// Format the response
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"id":           quiz.ID,
			"quiz":         quiz.Name,
			"question":     questionResults,
			"is_submitted": isSubmitted,
			"is_ended":     isEnded,
		},
		"meta": map[string]interface{}{
			"image_path": metaData.ImagePath,
			"video_path": metaData.VideoPath,
		},
		"status": map[string]interface{}{
			"code":    0,
			"message": "Success",
		},
	}

	return e.JSON(http.StatusOK, response)
}
