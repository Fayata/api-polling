package controllers

import (
	"api-polling/system/database"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func AddPoll(c echo.Context) error {
	// data request user
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"data": nil,
			"status": map[string]interface{}{
				"message": "user_id tidak valid",
				"code":    400,
				"type":    "error",
			},
		})
	}
	pollID, err := strconv.Atoi(c.Param("poll_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"data": nil,
			"status": map[string]interface{}{
				"message": "poll_id tidak valid",
				"code":    400,
				"type":    "error",
			},
		})
	}
	choiceID, err := strconv.Atoi(c.Param("choice_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"data": nil,
			"status": map[string]interface{}{
				"message": "choice_id tidak valid",
				"code":    400,
				"type":    "error",
			},
		})
	}
	// Validasi data

	db, err := database.Conn()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"data": nil,
			"status": map[string]interface{}{
				"message": "Gagal terhubung ke database",
				"code":    500,
				"type":    "error",
			},
		})
	}
	defer db.Close()

	query := "INSERT INTO user_choice (user_id, poll_id, choice_id) VALUES (?, ?, ?)"
	_, err = db.Exec(query, userID, pollID, choiceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"data": nil,
			"status": map[string]interface{}{
				"message": "Gagal menyimpan polling",
				"code":    500,
				"type":    "error",
			},
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": map[string]interface{}{
			"polling_id": pollID,
		},
		"status": map[string]interface{}{
			"message": "submit berhasil",
			"code":    0,
			"type":    "success",
		},
	})
}
