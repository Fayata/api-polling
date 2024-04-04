package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func AllResult(e echo.Context) error {
	var AllResults []*models.PollingResult

	db, err := database.Conn()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}
	defer db.Close()

	// Subquery untuk menghitung jumlah user yang memberikan suara
	rows, err := db.Query(`
        SELECT 
            polling.title,
            SUM(CASE WHEN result.vote = 1 THEN 1 ELSE 0 END) as item1_count,
            SUM(CASE WHEN result.vote = 2 THEN 1 ELSE 0 END) as item2_count,
            COUNT(DISTINCT result.user_id) as total_participants
        FROM polling
        LEFT JOIN result ON polling.poll_id = result.poll_id
        LEFT JOIN user ON result.user_id = user.user_id
        GROUP BY polling.poll_id
    `)
	if err != nil {
		log.Println("Failed to execute query:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var result models.PollingResult
		if err := rows.Scan(&result.Title, &result.Item1Count, &result.Item2Count, &result.TotalParticipants); err != nil {
			log.Println("Failed to scan row:", err)
			return err
		}

		result.Item1Percentage = float64(result.Item1Count) / float64(result.TotalParticipants) * 100
		result.Item2Percentage = float64(result.Item2Count) / float64(result.TotalParticipants) * 100

		AllResults = append(AllResults, &result)
	}
	return e.JSON(http.StatusOK, AllResults)
}
