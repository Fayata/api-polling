package main

import (
	"api-polling/application/config/app"
	"api-polling/application/config/routes"
	"fmt"
)

func main() {
	webServer := routes.AppRoute()
	webServerConfig := fmt.Sprintf("%v:%v", app.Load.WebServer.Host, app.Load.WebServer.Port)
	webServer.Logger.Fatal(webServer.Start(webServerConfig))
}