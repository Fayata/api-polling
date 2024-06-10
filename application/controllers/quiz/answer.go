package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"net/http"

	"github.com/labstack/echo"
)

func Answer(e echo.Context) error {

	var quizOption models.QuizQuestionChoice
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid option_id")
	}
	if err := db.First(&quizOption, quizOption.QuestionID).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid option_id")
	}
	var quiz models.Quiz
	err = db.First(&quiz.ID).Error
	if err != nil{
		return echo.NewHTTPError(http.StatusNotFound, "Quiz tidak ditemukan")
	}

	var quizQ models.QuizQuestion
	err = db.First(&quizOption.ID).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Option tidak ditemukan")
	}

	userAnswer := models.UserAnswer{
		QuestionID: uint(quizQ.ID),
		UserID:     int(e.Get("user_id").(int)),
		ChoiceID:   uint(quizOption.ID),
		QuizID:     uint(quiz.ID),
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
