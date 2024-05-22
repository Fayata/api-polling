package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"api-polling/application/models"
	"api-polling/system/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/gorm"
)

var jwtSecret = []byte(getEnv("JWT_SECRET", "default_secret"))

func getEnv(key, defaultValue string) string {
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
		if err != nil || !token.Valid {
			log.Println("Token tidak valid:", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Klaim token tidak valid")
			return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			log.Println("Format klaim user_id tidak valid")
			return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
		}

		// Pengecekan keberadaan user di database
		var user models.User
		if err := database.GetDB().First(&user, uint(userID)).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return echo.NewHTTPError(http.StatusUnauthorized, "User tidak ditemukan")
			}
			log.Println("Error saat mengambil data user:", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Terjadi kesalahan")
		}

		e.Set("user_id", int(userID))

		return next(e)
	}
}
