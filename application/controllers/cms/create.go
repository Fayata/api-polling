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

	query := "INSERT INTO polling (title) VALUES (?)"
	query2 := "INSERT INTO poll_choices (option, poll_id) VALUES (?, ?)"
	_, err = db.Exec(query, query2, newPoll.Title, newPoll.Choices)
	if err != nil {
		log.Println("Gagal membuat polling baru:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal membuat polling baru")
	}

	return e.JSON(http.StatusCreated, newPoll)
}
