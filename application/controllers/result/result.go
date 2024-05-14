package controllers

import (
	"api-polling/application/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func Result(e echo.Context) error {
	pollID, err := strconv.Atoi(e.Param("poll_id"))
	if err != nil {
		log.Println("Gagal mengkonversi poll_id:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "poll_id tidak valid")
	}

	// Mengambil hasil polling berdasarkan ID polling menggunakan model
	pollingResults, err := models.GetPollingResultsByID(pollID)
	if err != nil {
		log.Println("Gagal mengambil hasil polling:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil hasil polling")
	}

	return e.JSON(http.StatusOK, pollingResults)
}
