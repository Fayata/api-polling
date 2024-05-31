package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"net/http"

	"github.com/labstack/echo"
)



func Answer(e echo.Context) error {

	type AddQuizAnswerRequest struct {
		OptionID       uint `json:"option_id"`
		QuestionNumber uint `json:"question_number"`
	}
    var req AddQuizAnswerRequest
    if err := e.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    var quizOption models.QuizOption
    if err := database.GetDB().First(&quizOption, req.OptionID).Error; err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid option_id")
    }

    var quiz models.Quiz
    if err := database.GetDB().First(&quiz, quizOption.QuizID).Error; err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil data quiz")
    }

    userAnswer := models.UserQuizAnswer{
        OptionID: req.OptionID,
        UserID:   uint(e.Get("user_id").(int)),
        QuizID:   quizOption.QuizID,
    }

    if err := database.GetDB().Create(&userAnswer).Error; err != nil {
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
                "current": req.QuestionNumber,
            },
        },
        "status": map[string]interface{}{
            "code":    0,
            "message": "Submit berhasil",
        },
    }

    return e.JSON(http.StatusOK, response)
}
