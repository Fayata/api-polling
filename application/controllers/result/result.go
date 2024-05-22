package controllers

import (
	"api-polling/application/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func Result(e echo.Context) error {
	pollIDStr := e.Param("poll_id")
	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		log.Println("Gagal mengkonversi poll_id:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "poll_id tidak valid")
	}

	// Mengambil hasil polling berdasarkan ID polling menggunakan model
	pollingResults, err := models.GetPollingResultsByID(int(pollID)) // ubah jadi uint
	if err != nil {
		log.Println("Gagal mengambil hasil polling:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil hasil polling")
	}

	return e.JSON(http.StatusOK, pollingResults)
}
