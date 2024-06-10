package routes

import (
	controllers_cms "api-polling/application/controllers/cms"
	controllers_polling "api-polling/application/controllers/polling"
	controllers_quiz "api-polling/application/controllers/quiz"
	controllers_result "api-polling/application/controllers/result"
	middleware "api-polling/application/middleware"

	"github.com/labstack/echo"
)

func AppRoute() *echo.Echo {
	e := echo.New()
	middleware.SetCORS(e)

	// e.POST("api/v1/login", controllers_polling.Login)
	//Admin routes
	e.PUT("api/v1/cms/polling/:id", controllers_cms.Update)
	e.DELETE("api/v1/cms/polling/:id", controllers_cms.Delete)
	e.POST("api/v1/cms/polling", controllers_cms.Create)
	//Users Polling
	e.GET("api/v1/polling/all", controllers_polling.AllList, middleware.JWTMiddleware)
	e.GET("api/v1/polling/:id", controllers_polling.ByID)//, middleware.JWTMiddleware)
	e.POST("api/v1/polling/:id", controllers_polling.AddPoll)//, middleware.JWTMiddleware)
	e.GET("api/v1/polling/:poll_id/leaderboard", controllers_result.Result, middleware.JWTMiddleware)

	//Users Quiz
	e.GET("api/v1/quiz/:id", controllers_quiz.GetQuizByID, middleware.JWTMiddleware)
	e.POST("api/v1/quiz/:id", controllers_quiz.Answer, middleware.JWTMiddleware)

	return e
}
