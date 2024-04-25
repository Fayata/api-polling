package controllers

import (
	"api-polling/application/models"
	"net/http"
	"strconv"
	"github.com/labstack/echo"
)

// AddPoll menangani penambahan polling baru
func AddPoll(e echo.Context) error {
	// Parse data dari request
	userID, _ := strconv.Atoi(e.Param("user_id"))
	pollID, _ := strconv.Atoi(e.Param("poll_id"))
	choiceID, _ := strconv.Atoi(e.Param("choice_id"))

	// Membuat instance polling
	polling := models.Polling{}

	// Menambahkan polling
	if err := polling.AddPoll(userID, choiceID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal melakukan polling")
	}

	return e.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{
			"user_id":    userID,
			"poll_id":    pollID,
			"choice_id":  choiceID,
		},
		"status": map[string]interface{}{
			"message": "Polling berhasil",
			"code":    http.StatusOK,
			"type":    "success",
		},
	})
}
