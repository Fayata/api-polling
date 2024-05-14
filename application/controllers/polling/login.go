package controllers

import (
	"api-polling/application/models"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// getEnv mengambil nilai environment variable atau nilai default jika tidak ada
func getEnv(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

func Login(e echo.Context) error {
	request := models.User{}
	err := e.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request")
	}

	user, err := request.Login(request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal melakukan login")
	}

	//Signature token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["id"] = user.ID

	tokenString, err := token.SignedString([]byte(getEnv("jwt_secret", "default_secret")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal menyimpan token")
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt_token"
	cookie.Value = tokenString
	cookie.Path = "/"
	e.SetCookie(cookie)

	response := map[string]interface{}{
		"access_token": tokenString,
	}

	return e.JSON(http.StatusOK, response)
}
