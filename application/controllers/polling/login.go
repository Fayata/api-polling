package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"database/sql"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(e echo.Context) error {
	request := LoginRequest{}
	err := e.Bind(&request)
	if err != nil {
		log.Printf("error binding: %v\n ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request")
	}

	// Membuat koneksi ke database
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal terhubung ke database")
	}
	defer db.Close()

	//query untuk memangil email dan password
	query := "SELECT user_id, email, token FROM user WHERE email = ? AND password = ?"
	row := db.QueryRow(query, request.Email, request.Password)

	var user models.User

	err = row.Scan(&user.User_id, &user.Email, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			// Pengguna tidak ditemukan dalam database
			log.Println("Pengguna dengan email tersebut tidak ditemukan")
			return echo.NewHTTPError(http.StatusUnauthorized, "Email atau kata sandi salah")
		}
		log.Println("Gagal melakukan query:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil data user")
	}

	//buat token
	token := jwt.New(jwt.SigningMethodHS256)

	//Menambah klaim token
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.User_id

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println("Gagal menyimpan token:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menyimpan token")
	}
	response := map[string]interface{}{
		"access_token": tokenString,
		"user_id":      user.User_id,
	}

	return e.JSON(http.StatusOK, response)
}
