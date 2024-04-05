package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func AllResult(e echo.Context) error {
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	// Query untuk mendapatkan hasil poling untuk setiap item dan total partisipant
	query := `
		SELECT 
			p.title,
			ROUND(SUM(CASE WHEN r.vote = 1 THEN 1 ELSE 0 END) / COUNT(r.user_id) * 100, 2) AS item1_percentage,
			ROUND(SUM(CASE WHEN r.vote = 2 THEN 1 ELSE 0 END) / COUNT(r.user_id) * 100, 2) AS item2_percentage,
			ROUND(SUM(CASE WHEN r.vote = 3 THEN 1 ELSE 0 END) / COUNT(r.user_id) * 100, 2) AS item3_percentage,
			ROUND(SUM(CASE WHEN r.vote = 4 THEN 1 ELSE 0 END) / COUNT(r.user_id) * 100, 2) AS item4_percentage,
			ROUND(SUM(CASE WHEN r.vote = 5 THEN 1 ELSE 0 END) / COUNT(r.user_id) * 100, 2) AS item5_percentage,
			COUNT(r.user_id) AS total_participants
		FROM result r
		INNER JOIN polling p ON r.poll_id = p.poll_id
		GROUP BY r.poll_id
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Gagal melakukan query:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil data hasil polling")
	}
	defer rows.Close()

	var pollingResults []*models.PollingResult

	for rows.Next() {
		var pollingResult models.PollingResult
		if err := rows.Scan(&pollingResult.Title, &pollingResult.Item1Percentage, &pollingResult.Item2Percentage, &pollingResult.Item3Percentage, &pollingResult.Item4Percentage, &pollingResult.Item5Percentage, &pollingResult.TotalParticipants); err != nil {
			log.Println("Gagal memindai baris:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memindai hasil polling")
		}
		pollingResults = append(pollingResults, &pollingResult)
	}

	return e.JSON(http.StatusOK, pollingResults)
}
