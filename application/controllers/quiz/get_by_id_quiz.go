package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetQuizByID(e echo.Context) error {
    id, err := strconv.Atoi(e.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
    }

    var quiz models.Quiz
    var question models.QuizQuestion

    // Get quiz by ID
    err = database.GetDB().Preload("Options").First(&quiz, id).Error
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Quiz tidak ditemukan")
    }
	
    userID := e.Get("user_id").(int) 

    // Check if quiz is submitted and ended
    isSubmitted, isEnded, err := quiz.CheckQuizStatus(e, uint(userID)) 
    if err != nil {
        // Tangani kesalahan jika ada
        log.Println("Gagal memeriksa status quiz:", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memeriksa status quiz")
    }

    currentQuestion, err := quiz.GetQuizPosition() 
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mendapatkan posisi quiz")
    }

    // Format response sesuai dengan struktur yang diinginkan
    response := map[string]interface{}{
        "data": map[string]interface{}{
            "id":         quiz.ID,
            "title":      quiz.Name,
            "question":   quiz.Question,
            "option": map[string]interface{}{
                //"type": 
                "data": map[string]interface{}{
                    "id": question.ID,
                    "label":question.QuestionText,
                    "quiz_id":question.QuizID,
                    "image_url":question.QuestionImage,
                },
            },
            "banner": map[string]interface{}{
                "type": quiz.IsActive,
                "url":  question.QuestionImage,
            },
            "is_submitted": isSubmitted,
            "is_ended":     isEnded,
        },
        "meta": map[string]interface{}{
            "questions": map[string]interface{}{
                "total":   quiz.TotalQuestion,
                "current": currentQuestion,
            },
        },
        "status": map[string]interface{}{
            "code":    0,
            "message": "Success",
        },
    }

    return e.JSON(http.StatusOK, response)
}
