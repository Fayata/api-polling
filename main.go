package main

import(
	"api-polling/controllers"
	
	"github.com/labstack/echo"
)

func main(){
	e := echo.New()

	e.GET("api/v1/polling/:id", controllers.AllList)

	e.Logger.Fatal(e.Start(":9000"))
}