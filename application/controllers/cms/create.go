package controllers

import (
    "api-polling/application/models"
    "github.com/labstack/echo"
    "net/http"
)

func Create(e echo.Context) error {
    var newPoll models.Poll
    if err := e.Bind(&newPoll); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Gagal melakukan binding data")
    }

    err := newPoll.Create()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal membuat polling baru")
    }

    return e.JSON(http.StatusCreated, newPoll)
}