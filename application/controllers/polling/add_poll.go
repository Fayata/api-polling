package controllers

import (
    "api-polling/application/models"
    "net/http"
    "strconv"

    "github.com/labstack/echo"
)

type AddPollRequest struct {
    ChoiceID uint `json:"choice_id"`
}

func AddPoll(e echo.Context) error {
    var req AddPollRequest
    if err := e.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    userID := e.Get("user_id").(int)
    userChoice := models.UserChoice{
        ChoiceID: req.ChoiceID, 
        UserID:   int(userID),
    }

    pollIDStr := e.Param("id")
    pollID, err := strconv.Atoi(pollIDStr)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid poll_id")
    }
    userChoice.PollID = int(pollID)
	

    if err := userChoice.AddPoll(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    response := map[string]interface{}{
        "data": 0,
        "meta": 0,
        "status": map[string]interface{}{
            "message": "submit berhasil",
            "code":    0,
            "type":    "success",
        },
    }

    return e.JSON(http.StatusOK, response)
}
