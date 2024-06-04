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

    userChoice := models.User_Answer{
        Choice_id: req.Choice_id, 
        User_Id:   req.User_Id,
    }

    pollIDStr := e.Param("id")
    pollID, err := strconv.Atoi(pollIDStr)
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid poll_id")
    }
    userChoice.Poll_Id = int(pollID)
	

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
