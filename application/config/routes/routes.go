package routes

import (
	controllers_polling "api-polling/application/controllers/polling"
	controllers_cms "api-polling/application/controllers/cms"

	"github.com/labstack/echo"
)

func AppRoute() *echo.Echo {
	e := echo.New()

	//Admin routes
	e.PUT("api/v1/cms/polling/:id", controllers_cms.Update)
	e.DELETE("api/v1/cms/polling/:id", controllers_cms.Delete)
	e.POST("api/v1/cms/polling",controllers_cms.Create)
	//Users
	e.GET("api/v1/polling", controllers_polling.AllList)
	e.GET("api/v1/polling/:id", controllers_polling.ByID)
	e.POST("api/v1/polling/add", controllers_polling.AddPoll)
	e.GET("api/v1/polling/leaderboard/:poll_id", controllers_polling.Result)
	return e
}
