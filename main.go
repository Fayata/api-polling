package main

import(
	"api-polling/controllers"
	
	"github.com/labstack/echo"
)

func main(){
	e := echo.New()

	e.GET("api/v1/polling", controllers.AllList)
	// e.POST("api/v1/polling",controllers.Create)
	e.GET("api/v1/polling/result", controllers.AllResult)

	e.Logger.Fatal(e.Start(":9000"))
}