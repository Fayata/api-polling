package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func Result(e echo.Context) error {
	pollID, err := strconv.Atoi(e.Param("poll_id"))
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

	// Query untuk mendapatkan jumlah pemilih untuk setiap pilihan polling
	query := `
		SELECT pc.id, pc.poll_id, pc.title, pc.option, COUNT(uc.user_id) as vote_count
		FROM poll_choices pc
		LEFT JOIN user_choice uc ON pc.id = uc.choice_id
		WHERE pc.poll_id = ?
		GROUP BY pc.id, pc.poll_id,pc.title, pc.option
	`
	rows, err := db.Query(query, pollID)
	if err != nil {
		log.Println("Gagal melakukan query:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil hasil polling")
	}
	defer rows.Close()

	var pollingResults []map[string]interface{}

	for rows.Next() {
	
		var choice models.PollChoices
		var voteCount int
		if err := rows.Scan(&choice.ID, &choice.Poll_id, &choice.Title, &choice.Option, &voteCount); err != nil {
			log.Println("Gagal memindai baris:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memindai data hasil polling")
		}
		totalParticipants := getTotalParticipants(pollID) // Fungsi ini harus Anda definisikan untuk menghitung total peserta polling
		percentage := float64(voteCount) / float64(totalParticipants) * 100

		result := map[string]interface{}{
			"poll_choice": map[string]interface{}{
				"id":      choice.ID,
				"poll_id": choice.Poll_id,
				"option":  choice.Option,
				"title": choice.Title,
			},
			"percentage": percentage,
		}
		pollingResults = append(pollingResults, result)
	}

	return e.JSON(http.StatusOK, pollingResults)
}

func getTotalParticipants(pollID int) int {
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return 0
	}
	defer db.Close()

	var totalParticipants int
	err = db.QueryRow("SELECT COUNT(*) FROM poll_choices WHERE poll_id =?", pollID).Scan(&totalParticipants)
	if err != nil {
		log.Println("Gagal menjalankan query:", err)
		return 0
	}
	return totalParticipants
}
