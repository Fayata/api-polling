package controllers

import (
    "api-polling/application/models"
    "github.com/labstack/echo"
    "net/http"
)

func AllList(e echo.Context) error {
    var userPolling models.Polling

    if err := e.Bind(&userPolling); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    polls, err := userPolling.GetAll()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mendapatkan semua polling")
    }

    response := map[string]interface{}{
        "data": polls,
    }

    return e.JSON(http.StatusOK, response)
}
