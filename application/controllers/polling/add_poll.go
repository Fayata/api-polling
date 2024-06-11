package controllers

import (
	"api-polling/application/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func AddPoll(e echo.Context) error {
	var req models.User_Answer
	if err := e.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}


	userID := e.Get("user_id").(int)

	userChoice := models.User_Answer{
		Choice_id: req.Choice_id,
		User_Id:   userID,
	}

	pollIDStr := e.Param("id")
	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid poll_id")
	}
	userChoice.Poll_Id = int(pollID)

	// Panggil fungsi AddPoll dari model
	err = userChoice.AddPoll()

	// Buat respons dasar
	response := map[string]interface{}{
		"data": 0,
		"meta": 0,
		"status": map[string]interface{}{
			"message": "Submit Berhasil",
			"code":    0,
			"type":    "success",
		},
	}

	if err != nil {
		response["status"].(map[string]interface{})["code"] = 1 
		response["status"].(map[string]interface{})["type"] = "error"
		response["status"].(map[string]interface{})["message"] = err.Error()
		if err.Error() == "user has already polled" {
			return e.JSON(http.StatusBadRequest, response)
		} else {
			return e.JSON(http.StatusInternalServerError, response)
		}
	}

	return e.JSON(http.StatusOK, response)
}
