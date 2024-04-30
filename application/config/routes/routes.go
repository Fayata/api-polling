package routes

import (
	controllers_cms "api-polling/application/controllers/cms"
	controllers_polling "api-polling/application/controllers/polling"
	controllers_result "api-polling/application/controllers/result"
	middleware "api-polling/application/middleware"

	"github.com/labstack/echo"
)

func AppRoute() *echo.Echo {
	e := echo.New()

	e.POST("api/v1/polling/login", controllers_polling.Login)
	//Admin routes
	e.PUT("api/v1/cms/polling/:id", controllers_cms.Update)
	e.DELETE("api/v1/cms/polling/:id", controllers_cms.Delete)
	e.POST("api/v1/cms/polling", controllers_cms.Create)
	//Users
	e.GET("api/v1/polling/all", controllers_polling.AllList, middleware.JWTMiddleware)
	e.GET("api/v1/polling/:id", controllers_polling.ByID, middleware.JWTMiddleware)
	e.POST("api/v1/polling/:id", controllers_polling.AddPoll, middleware.JWTMiddleware)
	e.GET("api/v1/polling/leaderboard/:poll_id", controllers_result.Result, middleware.JWTMiddleware)
	return e
}
