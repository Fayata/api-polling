package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func Update(e echo.Context)error{
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil{
		log.Println("Failed to convert ID: ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	var updatePoll models.Polling
	if err := e.Bind(&updatePoll); err != nil{
		log.Println("Failed bind:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	db, err := database.Conn()
	if err != nil{
		log.Println("Failed Connect database", err)
		return err
	}
	defer db.Close()

	u, err := db.Exec("UPDATE polling SET title = ?, item1 = ?, item2 = ? WHERE poll_id = ?", updatePoll.Title, updatePoll.Item1, updatePoll.Item2, id)
	if err != nil{
		log.Println("Failed query:", err)
		return err
	}

	rowsAffected,_ := u.RowsAffected()
	if rowsAffected == 0{
		return echo.NewHTTPError(http.StatusNotFound, "Polling not found")
	}

	return e.JSON(http.StatusOK, echo.Map{
		"message": "Updsate successfully",
	})
}