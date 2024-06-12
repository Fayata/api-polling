package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
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

func GetQuizByID(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	var quiz models.Quiz

	// Get quiz by ID
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}
	err = db.First(&quiz, id).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Quiz not found")
	}

	userID := e.Get("user_id").(int)

	// Check if quiz is submitted and ended
	isSubmitted, isEnded, err := quiz.CheckQuizStatus(e, uint(userID))
	if err != nil {
		log.Println("Failed to check quiz status:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	// Fetch quiz questions with their Quiz
	questions, err := models.GetQuestionByQuizId(id)
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
		choiceDataArr, err := models.GetChoiceByQuestionId(question.ID)
		if err != nil {
			log.Println("Failed to get question choices:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		for i, choice := range choiceDataArr {
			choiceDataArr[i].QuestionType = models.GetChoiceType(choice.ChoiceImage)
			questionResult.Type = models.GetChoiceType(choice.ChoiceImage)
		}
		questionResult.Order.Total = len(questions)
		questionResult.Order.Current = i + 1
		questionResult.Choices = choiceDataArr
		// questionResult.Type = models.GetChoiceType(choice.ChoiceImage)

		// Get banner info (using GetQuestionType here, assuming it returns banner details)
		questionResult.Banner.ShowBanner = models.GetQuestionType(question.QuestionImage)
		questionResult.Banner.URL = question.QuestionImage

		questionResults = append(questionResults, questionResult)
	}

	// Format the response
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"id":           quiz.ID,
			"title":        quiz.Name,
			"question":     questionResults,
			"is_submitted": isSubmitted,
			"is_ended":     isEnded,
		},
		"meta": map[string]interface{}{
			"image_path": database.GetAppConfig().ImagePath,
			"video_path": database.GetAppConfig().VideoPath,
		},
		"status": map[string]interface{}{
			"code":    0,
			"message": "Success",
		},
	}

	return e.JSON(http.StatusOK, response)
}
