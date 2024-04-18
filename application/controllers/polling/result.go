package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func Result(c echo.Context) error {
	pollID, err := strconv.Atoi(c.Param("poll_id"))
	if err != nil {
		log.Println("Gagal mengkonversi poll_id:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "poll_id tidak valid")
	}
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	// Query hasil poling untuk setiap item dan total participant
	query := `
		SELECT p.title,
			ROUND(SUM(CASE WHEN r.vote = 1 THEN 1 ELSE 0 END) / COUNT(DISTINCT r.user_id) * 100, 2) AS item1_percentage,
			ROUND(SUM(CASE WHEN r.vote = 2 THEN 1 ELSE 0 END) / COUNT(DISTINCT r.user_id) * 100, 2) AS item2_percentage,
			ROUND(SUM(CASE WHEN r.vote = 3 THEN 1 ELSE 0 END) / COUNT(DISTINCT r.user_id) * 100, 2) AS item3_percentage,
			ROUND(SUM(CASE WHEN r.vote = 4 THEN 1 ELSE 0 END) / COUNT(DISTINCT r.user_id) * 100, 2) AS item4_percentage,
			ROUND(SUM(CASE WHEN r.vote = 5 THEN 1 ELSE 0 END) / COUNT(DISTINCT r.user_id) * 100, 2) AS item5_percentage,
			COUNT(DISTINCT r.user_id) AS total_participants
		FROM result r
		INNER JOIN polling p ON r.poll_id = p.poll_id
		WHERE r.poll_id = ?
		GROUP BY r.poll_id
	`
	var pollingResult models.PollingResult

	err = db.QueryRow(query, pollID).Scan(&pollingResult.Title, &pollingResult.Item1Percentage, &pollingResult.Item2Percentage, &pollingResult.Item3Percentage, &pollingResult.Item4Percentage, &pollingResult.Item5Percentage, &pollingResult.TotalParticipants)
	if err != nil {
		log.Println("Gagal mendapatkan hasil polling:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mendapatkan hasil polling")
	}

	return c.JSON(http.StatusOK, pollingResult)
}
