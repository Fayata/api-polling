// middlewares/auth.go
package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)
var jwtSecret = []byte(getEnv("JWT_SECRET", "default_secret"))

func getEnv(key, defaultValue string) string {
    val, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }
    return val
}
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(e echo.Context) error {
        // Ambil token dari header Authorization
        authHeader := e.Request().Header.Get("Authorization")
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

        // Parse dan verifikasi token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            log.Println("Token tidak valid:", err)
            return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
        }

        // Ambil klaim dari token
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            log.Println("Klaim token tidak valid")
            return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
        }

        // Set klaim token ke dalam konteks
        e.Set("user_id", claims["user_id"])

        // Lanjutkan ke handler berikutnya
        return next(e)
    }
}