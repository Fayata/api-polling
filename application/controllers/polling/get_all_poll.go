package controllers

import (
    "api-polling/application/models"
    "github.com/labstack/echo"
    "net/http"
)

func AllList(e echo.Context) error {
   
    // Membuat slice untuk menyimpan semua polling
    polling := models.Polling{}

    //getallpoll
    polls, err := polling.GetAll()
    if err!= nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    // Mengembalikan data polling dalam format JSON
    return e.JSON(http.StatusOK, polls)
}
