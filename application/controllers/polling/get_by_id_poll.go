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
		return echo.NewHTTPError(http.StatusBadRequest, "Poll_id tidak valid")
	}

	db := database.GetDB()

	// Mengambil semua data yang diperlukan
	var polling models.Poll
	if err := db.Find(&polling, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Polling tidak ditemukan")
	}

	var pollChoices []models.Poll_Choices
	if err := db.Where("poll_id = ?", id).Find(&pollChoices).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil pilihan polling")
	}

	var pollResults []models.Poll_Result
	if err := db.Where("poll_id = ?", id).Find(&pollResults).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil hasil polling")
	}

	userIDInterface := e.Get("user_id")
	var userAnswers []models.User_Answer
	if userIDInterface != nil {
		userID, ok := userIDInterface.(int)
		if ok {
			if err := db.Where("user_id = ? AND poll_id = ?", userID, id).Find(&userAnswers).Error; err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil jawaban pengguna")
			}
		}
	}

	// Check if poll is submitted and ended (sesuaikan logika ini)
	isSubmitted, isEnded := false, false
	for _, ua := range userAnswers {
		if ua.Poll_Id == id {
			isSubmitted = true
			break
		}
	}

	if isSubmitted && len(userAnswers) == len(pollChoices) {
		isEnded = 
		true
	}
	

	// Format pilihan (choices) untuk respons
	formattedChoices := make([]map[string]interface{}, len(pollChoices))
	for i, choice := range pollChoices {
		isSelected := false
		for _, ua := range userAnswers {
			if ua.Choice_id == choice.ID {
				isSelected = true
				break
			}
		}

		// Dapatkan vote count dari pollResults
		voteCount := 0
		for _, pr := range pollResults {
			if pr.Choice_id == choice.ID {
				voteCount++
			}
		}

		formattedChoices[i] = map[string]interface{}{
			"id":          choice.ID,
			"label":       choice.Choice_text,
			"image_url":   choice.Choice_image,
			"value":       voteCount,
			"is_selected": isSelected,
		}
	}


	response := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         polling.ID,
			"title":      polling.Title,
			"question":   polling.Question,
			"option": map[string]interface{}{
				"type": polling.Option_type,
				"data": formattedChoices, 
			},
			"banner": map[string]interface{}{
				"type": polling.Question_text,
				"url":  polling.Question_image,
			},
			"is_submitted": isSubmitted,
			"is_ended":     isEnded,
		},
		"meta": map[string]interface{}{
			"questions": map[string]interface{}{
				"total":   0,
				"current": 0,
			},
		},
		"status": map[string]interface{}{
			"code":    0,
			"message": "Success",
		},
	}

	return e.JSON(http.StatusOK, response)
}
