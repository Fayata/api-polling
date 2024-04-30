package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var jwtSecret = []byte(getEnv("jwt_secret", "default_secret"))

func getEnv(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
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
		e.Set("user_id", int(userID))

		return next(e)
	}
}
