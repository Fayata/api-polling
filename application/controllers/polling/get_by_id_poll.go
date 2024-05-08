package controllers

import (
    "api-polling/application/models"
    "github.com/labstack/echo"
    "net/http"
    "strconv"
)

func ByID(e echo.Context) error {
   id, err := strconv.Atoi(e.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
    }

    var polling models.Polling

    // Get polling by ID
    err = polling.GetByID(id)
    if err != nil {
        return echo.NewHTTPError(http.StatusNotFound, "Polling tidak ditemukan")
    }

    return e.JSON(http.StatusOK, polling)
}