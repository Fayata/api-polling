package controllers

import (
	"api-polling/application/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func AddPoll(e echo.Context) error {
    // Parse data dari request
    userID, _ := strconv.Atoi(e.Param("user_id"))
    pollID, _ := strconv.Atoi(e.Param("poll_id"))
    choiceID, _ := strconv.Atoi(e.Param("choice_id"))


    tokenCookie, err := e.Cookie("jwt_token")
    if err != nil {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak ditemukan")
    }
    tokenString := tokenCookie.Value

    // Parse dan verifikasi token JWT
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("secret"), nil
    })
    if err != nil || !token.Valid {
        return echo.NewHTTPError(http.StatusUnauthorized, "Token tidak valid")
    }

    // Create polling instance
    polling := models.Polling{}

    // Add poll
    if err := polling.AddPoll(userID, choiceID); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Gagal melakukan polling")
    }

    return e.JSON(http.StatusCreated, map[string]interface{}{
        "data": map[string]interface{}{
            "user_id":  userID,
            "poll_id":  pollID,
            "choice_id": choiceID,
        },
        "status": map[string]interface{}{
            "message": "Polling berhasil",
            "code":    http.StatusOK,
            "type":    "success",
        },
    })
}
