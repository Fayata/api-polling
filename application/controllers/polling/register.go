package controllers

// import (
//     "api-polling/application/models"
//     "net/http"

//     "github.com/labstack/echo"
// )

// func Register(e echo.Context) error {
//     var user models.User
//     if err := e.Bind(&user); err != nil { // Bind to the address of 'user'
//         return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//     }

//     if err := user.Register()(); err != nil {
//         return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
//     }

//     return e.JSON(http.StatusCreated, user) // Return the created user with 201 Created status
// }
