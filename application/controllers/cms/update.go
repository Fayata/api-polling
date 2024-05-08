package controllers

import (
    "api-polling/application/models"
    "github.com/labstack/echo"
    "net/http"
    "strconv"
)

func Update(e echo.Context) error {
    id, err := strconv.Atoi(e.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
    }

    var updatePoll models.Polling
    if err := e.Bind(&updatePoll); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }

    err = updatePoll.Update(id)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengupdate data polling")
    }

    return e.JSON(http.StatusOK, echo.Map{
        "message": "Update successfully",
    })
}
