package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func Answer(e echo.Context) error {

	var req models.UserAnswer
	var quiz models.Quiz
	var quizOption models.QuizQuestionChoice
	if err := e.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid option_id")
	}


	userChoice := models.UserAnswer{
		ChoiceID: req.ChoiceID,
		UserID:   req.UserID,
	}

	quizIDStr := e.Param("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid poll_id")
	}
	userChoice.QuizID = int(quizID)

	err = userChoice.AddQuiz()
	if err != nil{
		return err
	}

	userAnswer := models.UserAnswer{
		QuestionID: int(e.Get("question_id").(int)),
		UserID:     int(e.Get("user_id").(int)),
		ChoiceID:   uint(quizOption.ID),
		QuizID:     int(quiz.ID),
	}

	if err := db.Create(&userAnswer).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menyimpan jawaban")
	}

	totalQuizzes, err := quiz.GetTotalQuizzes()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menghitung total quizzes")
	}

	response := map[string]interface{}{
		"data": "",
		"meta": map[string]interface{}{
			"questions": map[string]interface{}{
				"total":   totalQuizzes,
				"current": 0,
			},
		},
		"status": map[string]interface{}{
			"code":    0,
			"message": "Submit berhasil",
		},
	}

	return e.JSON(http.StatusOK, response)
}
