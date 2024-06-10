package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type TokenRplus struct {
	Vid       int    `json:"vid"`
	Token     string `json:"token"`
	Pl        string `json:"pl"`
	Device    string `json:"device_id"`
	LoginType string `json:"ltype"`
}

type TokenExisting struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func getEnv(key, defaultValue string) string {
	log.Printf("===> %v", os.Getenv("JWT_SECRET"))

	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

func SetCORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	var jwtSecret = []byte(getEnv("JWT_SECRET", "default_secret"))

	return func(e echo.Context) error {
		authHeader := e.Request().Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Token tidak ditemukan")
			return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak ditemukan")
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		log.Println("Token yang diterima:", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		// if err != nil || !token.Valid {
		// 	log.Println("Token tidak valid:", err)
		// 	return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak validi")
		// }

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Klaim token tidak valid")
			return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valied")
		}

		claimsByte, err := json.Marshal(claims)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed marshal")
		}
		var tokenRplus TokenRplus
		json.Unmarshal(claimsByte, &tokenRplus)

		// var user models.User
		// if err := database.GetDB().First(&user, tokenRplus.Vid).Error; err != nil {
		// 	if err == gorm.ErrRecordNotFound {
		// 		return echo.NewHTTPError(http.StatusUnauthorized, "User tidak ditemukan")
		// 	}
		// 	log.Println("Error saat mengambil data user:", err)
		// 	return echo.NewHTTPError(http.StatusInternalServerError, "Terjadi kesalahan")
		// }

		e.Set("user_id", tokenRplus.Vid)

		return next(e)
	}
}
