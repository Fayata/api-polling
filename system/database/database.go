package database

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type PollingDB struct {
    Db *gorm.DB
}

type QuizDB struct {
    Db *gorm.DB
}

type AppConfig struct {
    ImagePath string
    VideoPath string
}

func InitDB() (*PollingDB, *QuizDB) {
    loadEnv()

    pollingDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_POLLING_USER"),
        os.Getenv("DB_POLLING_PASSWORD"),
        os.Getenv("DB_POLLING_HOST"),
        os.Getenv("DB_POLLING_PORT"),
        os.Getenv("DB_POLLING_NAME"),
    )
    quizDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_QUIZ_USER"),
        os.Getenv("DB_QUIZ_PASSWORD"),
        os.Getenv("DB_QUIZ_HOST"),
        os.Getenv("DB_QUIZ_PORT"),
        os.Getenv("DB_QUIZ_NAME"),
    )

    pollingDBConnection, err := gorm.Open(mysql.Open(pollingDSN), &gorm.Config{})
    if err != nil {
        log.Fatal("Error connecting to polling database:", err)
    }
    quizDBConnection, err := gorm.Open(mysql.Open(quizDSN), &gorm.Config{})
    if err != nil {
        log.Fatal("Error connecting to quiz database:", err)
    }

    pollingDB := &PollingDB{Db: pollingDBConnection}
    quizDB := &QuizDB{Db: quizDBConnection}

    return pollingDB, quizDB
}

func loadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func (d *PollingDB) GetDb() *gorm.DB {
    return d.Db
}

func (d *QuizDB) GetDb() *gorm.DB {
    return d.Db
}

func Meta() (appConfig AppConfig) {
    loadEnv()
    appConfig.ImagePath = os.Getenv("APP_IMAGE_PATH")
    appConfig.VideoPath = os.Getenv("APP_VIDIO_PATH")
    return appConfig
}
