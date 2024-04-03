package controllers

import (
	"log"
	"net/http"
	"strconv"

	"api-polling/system/database"

	"github.com/labstack/echo"
)

func Delete(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		log.Println("Gagal mengonversi ID menjadi integer:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
	}

	// Membuat koneksi ke database
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	// Menghapus terlebih dahulu semua hasil yang terkait dengan polling
	_, err = db.Exec("DELETE FROM result WHERE poll_id=?", id)
	if err != nil {
		log.Println("Gagal menghapus hasil polling:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menghapus hasil polling")
	}

	// Kemudian hapus baris polling dari tabel polling
	_, err = db.Exec("DELETE FROM polling WHERE poll_id=?", id)
	if err != nil {
		log.Println("Gagal menghapus data polling:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menghapus data polling")
	}

	return e.NoContent(http.StatusOK)
}
