package routes

import (
	"api-polling/application/controllers/polling"
	"github.com/labstack/echo"
)

func AppRoute() *echo.Echo {
	e := echo.New()

	e.GET("api/v1/polling", controllers.AllList)
	e.GET("api/v1/polling/:id", controllers.ByID)
	e.PUT("api/v1/polling/:id", controllers.Update)
	e.DELETE("api/v1/polling/:id", controllers.Delete)
	e.POST("api/v1/polling",controllers.Create)
	e.POST("api/v1/polling/", controllers.AddPoll)
	e.GET("api/v1/polling/result", controllers.Result)
	return e
}
