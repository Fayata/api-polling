package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func AddPoll(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Println("Gagal mengkonversi user_id:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "user_id tidak valid")
	}
	pollID, err := strconv.Atoi(c.Param("poll_id"))
	if err != nil {
		log.Println("Gagal mengkonversi poll_id:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "poll_id tidak valid")
	}
	vote := c.Param("vote")

	var newPoll models.Result
	if err := c.Bind(&newPoll); err != nil {
		log.Println("Gagal melakukan binding data:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Gagal melakukan binding data")
	}

	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	queryInsert := "INSERT INTO result (vote, user_id, poll_id) VALUES (?, ?, ?)"
	_, err = db.Exec(queryInsert, vote, userID, pollID)
	if err != nil {
		log.Println("Gagal membuat hasil polling baru:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal membuat hasil polling baru")
	}

	return c.JSON(http.StatusCreated, "Polling berhasil ditambahkan")
}
