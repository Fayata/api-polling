package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func ByID(e echo.Context) error {
    id, err := strconv.Atoi(e.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
    }

    var polling models.Polling

    // Get polling by ID
    err = database.GetDB().Preload("Choices").Preload("Banner").First(&polling, id).Error
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Polling tidak ditemukan")
    }

    // Check if poll is submitted and ended
    isSubmitted, isEnded, err := polling.CheckPollStatus(e)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memeriksa status polling")
    }

    totalQuestions := 1
    currentQuestion := 1

    // Format response sesuai dengan struktur yang diinginkan
    response := map[string]interface{}{
        "data": map[string]interface{}{
            "id":         polling.ID,
            "title":      polling.Title,
            "question":   polling.Question,
            "option": map[string]interface{}{
                "type": "image",
                 "data": polling,
            },
            "banner": map[string]interface{}{
                "type": polling.Type,
                "url":  polling.URL,
            },
            "is_submitted": isSubmitted,
            "is_ended":     isEnded,
        },
        "meta": map[string]interface{}{
            "questions": map[string]interface{}{
                "total":   totalQuestions,
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
