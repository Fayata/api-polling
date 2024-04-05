package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func AllList(e echo.Context) error {
	var PollList []*models.PollingResponse

	db, err := database.Conn()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT poll_id, title, item1, item2, item3, item4, item5 FROM polling")
	if err != nil {
		log.Println("Failed to execute query:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var polling models.Polling
		if err := rows.Scan(&polling.Poll_id, &polling.Title, &polling.Item1, &polling.Item2, &polling.Item3, &polling.Item4, &polling.Item5); err != nil {
			log.Println("Failed to scan row:", err)
			return err
		}
		PollList = append(PollList, &models.PollingResponse{
			ID:    polling.Poll_id,
			Title: polling.Title,
			Items: []*models.PollingItem{
				{Value: polling.Item1},
				{Value: polling.Item2},
				{Value: polling.Item3},
				{Value: polling.Item4},
				{Value: polling.Item5},
			},
		})
	}
	return e.JSON(http.StatusOK, PollList)
}
