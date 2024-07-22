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
    var (
        isAutoMigrate = app.Load.Database.AutoMigrate
    )

    // Load .env file
    err := godotenv.Overload()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initialize database connection
    pollingDB, quizDB := database.InitDB()
    // defer pollingDB.Db.Close()
    // defer quizDB.Db.Close()

    // AutoMigrate models
    if isAutoMigrate {
        err := pollingDB.Db.AutoMigrate(
            &models.Poll{},
            &models.Poll_Choices{},
            &models.User_Answer{},
            &models.Poll_Result{},
        )
        if err != nil {
            log.Fatal("Error migrating polling database:", err)
        }

        err = quizDB.Db.AutoMigrate(
            &models.Quiz{},
            &models.QuizQuestion{},
            &models.QuizQuestionChoice{},
            &models.UserAnswer{},
        )
        if err != nil {
            log.Fatal("Error migrating quiz database:", err)
        }
    }
	webServer := routes.AppRoute()
	webServerConfig := fmt.Sprintf("%v:%v", app.Load.WebServer.Host, app.Load.WebServer.Port)
	webServer.Logger.Fatal(webServer.Start(webServerConfig))
}
