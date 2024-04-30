package controllers

import (
	"api-polling/application/models"
	"net/http"
	"github.com/labstack/echo"
)

// AddPoll menangani penambahan polling baru
func AddPoll(e echo.Context) error {
    var userChoice models.UserChoice
    if err := e.Bind(&userChoice); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    if err := userChoice.AddPoll(); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add poll")
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
