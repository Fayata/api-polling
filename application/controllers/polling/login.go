package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func Login(e echo.Context) error {
	// Membuat koneksi ke database
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	query := "SELECT user_id, name, email FROM users"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Gagal melakukan query:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil data pengguna")
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.User_id, &user.Name, &user.Email); err != nil {
			log.Println("Gagal memindai baris:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memindai data pengguna")
		}
		users = append(users, &user)
	}

	return e.JSON(http.StatusOK, users)
}
