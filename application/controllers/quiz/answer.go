package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"net/http"

	"github.com/labstack/echo"
)



func Answer(e echo.Context) error {


    var quizOption models.QuizQuestionChoice
    if err := database.GetDB().First(&quizOption, quizOption.QuestionID).Error; err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid option_id")
    }

    var quiz models.Quiz
    if err := database.GetDB().First(&quiz, quizOption.ID).Error; err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil data quiz")
    }
    var quizQ models.QuizQuestion
   

    userAnswer := models.UserAnswer{
        QuestionID: uint(quizQ.ID),
        UserID:   int(e.Get("user_id").(int)),
        ChoiceID: uint(quizOption.ID),
        QuizID:   uint(quiz.ID),
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
