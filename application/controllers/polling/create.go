package controllers

import (
	"log"
	"net/http"

	"api-polling/application/models"
	"api-polling/system/database"

	"github.com/labstack/echo"
)

func Create(e echo.Context) error {
	var newPoll models.Polling
	if err := e.Bind(&newPoll); err != nil {
		log.Println("Gagal melakukan binding data:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Gagal melakukan binding data")
	}
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	query := "INSERT INTO polling (title, item1, item2) VALUES (?, ?, ?)"
	_, err = db.Exec(query, newPoll.Title, newPoll.Item1, newPoll.Item2)
	if err != nil {
		log.Println("Gagal membuat polling baru:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal membuat polling baru")
	}

	return e.JSON(http.StatusCreated, newPoll)
}
