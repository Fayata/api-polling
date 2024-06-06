package controllers

import (
	"api-polling/application/entities/dmo/polling"
	"net/http"

	"github.com/labstack/echo"
)

func Create(e echo.Context) error {
	var newPoll polling.Poll
	if err := e.Bind(&newPoll); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Gagal melakukan binding data")
	}

	err := newPoll.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal membuat polling baru")
	}

	return e.JSON(http.StatusCreated, newPoll)
}
