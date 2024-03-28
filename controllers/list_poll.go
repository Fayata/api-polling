package controllers

import (
	"api-polling/models"
	"net/http"
	"api-polling/routes"
	"github.com/labstack/echo"
)

func AllList(e echo.Context) error {
	var PollList []*models.Polling

	db, err := routes.Conn()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT poll_id, title, item1, item2 FROM polling")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var polling models.Polling
		if err := rows.Scan(&polling.Poll_id, &polling.Title, &polling.Item1, &polling.Item2); err != nil {
			return err
		}
		PollList = append(PollList, &polling)
	}
	return e.JSON(http.StatusOK, PollList)

	

}
