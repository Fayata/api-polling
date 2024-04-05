package controllers

import (
	"log"
	"net/http"
	"strconv"
	"api-polling/application/models"
	"api-polling/system/database"

	"github.com/labstack/echo"
)

func ByID(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		log.Println("Failed to convert ID to integer:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	var polling models.Polling

	db, err := database.Conn()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT poll_id, title, item1, item2 FROM polling WHERE poll_id = ?", id).Scan(&polling.Poll_id, &polling.Title, &polling.Item1, &polling.Item2, &polling.Item3, &polling.Item4, &polling.Item5)
	if err != nil {
		log.Println("Failed to execute query or no rows found:", err)
		return echo.NewHTTPError(http.StatusNotFound, "Polling not found")
	}

	pollingResponse := &models.PollingResponse{
		ID:    polling.Poll_id,
		Title: polling.Title,
		Items: []*models.PollingItem{
			{Value: polling.Item1},
			{Value: polling.Item2},
			{Value: polling.Item3},
			{Value: polling.Item4},
			{Value: polling.Item5},
		},
	}

	return e.JSON(http.StatusOK, pollingResponse)
}
