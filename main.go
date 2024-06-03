package main

import (
	"fmt"
	"log"

	"api-polling/application/config/app"
	"api-polling/application/config/routes"
	"api-polling/application/models"

	"api-polling/system/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Overload()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	database.InitDB()

	// AutoMigrate models
	err = database.GetDB().AutoMigrate(
		&models.User{},
		&models.Polling{},
		&models.PollChoice{},
		&models.UserChoice{},
		&models.Quiz{},
		&models.QuizOption{},
		&models.UserQuizAnswer{},
	)
	if err != nil {
		log.Fatal("Error migrating database:", err)
	}

	webServer := routes.AppRoute()
	webServerConfig := fmt.Sprintf("%v:%v", app.Load.WebServer.Host, app.Load.WebServer.Port)
	webServer.Logger.Fatal(webServer.Start(webServerConfig))
}
