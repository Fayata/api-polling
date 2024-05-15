package controllers

import (
	"api-polling/application/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func AddPoll(e echo.Context) error {
	var userChoice models.UserChoice
	if err := e.Bind(&userChoice); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	userID := e.Get("user_id").(int)
	userChoice.User_id = userID

	pollID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid poll_id")
    }
    userChoice.Poll_id = pollID

	if err := userChoice.AddPoll(e); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Polling hanya bisa sekali yah")
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"id": userChoice.Choice_ID,
			"status": map[string]interface{}{
				"message": "submit berhasil",
				"code":    0,
				"type":    "success",
			},
		},
	}

	return e.JSON(http.StatusOK, response)
}
