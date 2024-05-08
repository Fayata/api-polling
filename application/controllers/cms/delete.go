package controllers

import (
    "api-polling/application/models"
    "github.com/labstack/echo"
    "net/http"
    "strconv"
)

func Delete(e echo.Context) error {
    id, err := strconv.Atoi(e.Param("id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
    }

    var poll models.Polling
    err = poll.Delete(id)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menghapus data polling")
    }

    return e.NoContent(http.StatusOK)
}
