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

	var polling models.Polling
	var user models.UserChoice 

	// Get polling by ID
	err = database.GetDB().Preload("Choices").First(&polling, id).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Polling tidak ditemukan")
	}

	// Check if poll is submitted and ended
	userId := user.ID
	isSubmitted, isEnded, err := polling.CheckPollStatus(e, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memeriksa status polling")
	}

	totalQuestions := 0
	currentQuestion := 0

	// Format pilihan (choices) untuk respons
	formattedChoices := make([]map[string]interface{}, len(polling.Choices))
	for i, choice := range polling.Choices {
		isSelected := false
		for _, selectedID := range polling.UserC {
			if choice.ID == selectedID.ID {
				isSelected = true
				break
			}
		}

		// Dapatkan vote count dari fungsi di model PollChoice
		voteCount, err := choice.GetVoteCount()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil jumlah suara")
		}

		formattedChoices[i] = map[string]interface{}{
			"id":          choice.ID,
			"label":       choice.Option,
			"image_url":   choice.ImageURL,
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
				"type": polling.Type,
				"data": formattedChoices, // Gunakan pilihan yang sudah diformat
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
