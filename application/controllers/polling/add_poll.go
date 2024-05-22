package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"net/http"
	"strconv"

	// "strconv"

	"github.com/labstack/echo"
)

func AddPoll(e echo.Context) error {
    var req models.AddPollRequest
    if err := e.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    var pollChoice models.PollChoice
    if err := database.GetDB().First(&pollChoice, req.OptionID).Error; err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid option_id")
    }

    var polling models.Polling
    if err := database.GetDB().First(&polling, pollChoice.PollID).Error; err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil data polling")
    }

    userChoice := models.UserChoice{
        ChoiceID: int(req.OptionID),
        UserID:   int(e.Get("user_id").(int)),
        PollID:   pollChoice.PollID,
    }

    if err := userChoice.AddPoll(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    totalPolls, err := polling.GetTotalPolls()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menghitung total polls")
    }

    currentQuestion, err := strconv.ParseUint(e.QueryParam("question_number"), 10, 64)
    if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "Invalid question_number")
}

    response := map[string]interface{}{
        "data": "", 
        "meta": map[string]interface{}{
            "questions": map[string]interface{}{
                "total":   totalPolls,
                "current": currentQuestion, 
            },
        },
        "status": map[string]interface{}{
            "code":    0,
            "message": "Submit berhasil",
        },
    }

    return e.JSON(http.StatusOK, response)
}
