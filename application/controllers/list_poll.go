package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func AllList(e echo.Context) error {
	var PollList []*models.Polling

	db, err := database.Conn()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT poll_id, title, item1, item2 FROM polling")
	if err != nil {
		log.Println("Failed to execute query:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var polling models.Polling
		if err := rows.Scan(&polling.Poll_id, &polling.Title, &polling.Item1, &polling.Item2); err != nil {
			log.Println("Failed to scan row:", err)
			return err
		}
		PollList = append(PollList, &polling)
	}
	return e.JSON(http.StatusOK, PollList)
}
