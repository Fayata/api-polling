package routes

import (
	"api-polling/application/controllers"

	"github.com/labstack/echo"
)

func AppRoute() *echo.Echo {
	e := echo.New()

	e.GET("api/v1/polling", controllers.AllList)
	// e.POST("api/v1/polling",controllers.Create)
	e.GET("api/v1/polling/result", controllers.AllResult)

	return e
}
