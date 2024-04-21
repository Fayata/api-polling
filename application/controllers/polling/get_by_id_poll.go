package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func ByID(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		log.Println("Gagal mengkonversi ID menjadi integer:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
	}

	var polling models.Polling

	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	// Ambil data polling berdasarkan poll_id
	err = db.QueryRow("SELECT poll_id, title FROM polling WHERE poll_id = ?", id).
		Scan(&polling.ID, &polling.Title)

	if err != nil {
		log.Println("Gagal menjalankan query atau tidak ada baris yang ditemukan:", err)
		return echo.NewHTTPError(http.StatusNotFound, "Polling tidak ditemukan")
	}

	// Ambil data pilihan polling berdasarkan poll_id
	rows, err := db.Query("SELECT id, option FROM poll_choices WHERE poll_id = ?", id)
	if err != nil {
		log.Println("Gagal menjalankan query:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil pilihan polling")
	}
	defer rows.Close()

	var choices []models.PollChoice

	// Iterasi melalui hasil query dan memasukkan ke slice choices
	for rows.Next() {
		var choice models.PollChoice
		if err := rows.Scan(&choice.ID, &choice.Option); err != nil {
			log.Println("Gagal membaca baris dari hasil query:", err)
			continue
		}
		choices = append(choices, choice)
	}

	// Jika ada kesalahan saat iterasi
	if err := rows.Err(); err != nil {
		log.Println("Gagal membaca semua baris dari hasil query:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal membaca pilihan polling")
	}

	polling.Choices = choices

	return e.JSON(http.StatusOK, polling)
}
