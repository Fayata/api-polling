package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"github.com/labstack/echo"
)

func AddPoll(e echo.Context, user_id int, poll_id int, vote int) error {
	var newPoll models.Result
	if err := e.Bind(&newPoll); err != nil {
		log.Println("Gagal melakukan binding data:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Gagal melakukan binding data")
	}

	// Memastikan user_id dan poll_id valid
	if user_id <= 0 || poll_id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id atau poll_id tidak valid")
	}

	// Memeriksa apakah vote berada dalam rentang yang valid
	if vote < 1 || vote > 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "vote tidak valid")
	}

	// Membuat koneksi ke database
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memulai transaksi")
	}
	defer tx.Rollback()

	// Memperbarui participant pada item yang dipilih
	queryUpdate := "UPDATE polling SET participant = participant + 1 WHERE poll_id = ? AND (item1 = ? OR item2 = ? OR item3 = ? OR item4 = ? OR item5 = ?)"
	_, err = tx.Exec(queryUpdate, poll_id, vote, vote, vote, vote, vote)
	if err != nil {
		log.Println("Gagal memperbarui participant:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memperbarui participant")
	}

	// Memasukkan hasil polling ke dalam database
	queryInsert := "INSERT INTO result (vote, participant, user_id, poll_id) VALUES (?, 1, ?, ?)"
	_, err = tx.Exec(queryInsert, vote, user_id, poll_id)
	if err != nil {
		log.Println("Gagal membuat hasil polling baru:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal membuat hasil polling baru")
	}

	// Commit transaksi
	err = tx.Commit()
	if err != nil {
		log.Println("Gagal commit transaksi:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal commit transaksi")
	}

	return e.JSON(http.StatusCreated, "Polling berhasil ditambahkan")
}
