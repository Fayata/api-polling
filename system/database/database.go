package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {}

func InitDB() *Database {
    return &Database{}
}


func loadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func (d *Database) DbPolling()(DBPolling *gorm.DB, err error){
    loadEnv()
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            os.Getenv("DB_POLLING_USER"),
            os.Getenv("DB_POLLING_PASSWORD"),
            os.Getenv("DB_POLLING_HOST"),
            os.Getenv("DB_POLLING_PORT"),
            os.Getenv("DB_POLLING_NAME"),
        )
        return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func (d *Database) DbQuiz()(DBPolling *gorm.DB, err error){
    loadEnv()
        qsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            os.Getenv("DB_QUIZ_USER"),
            os.Getenv("DB_QUIZ_PASSWORD"),
            os.Getenv("DB_QUIZ_HOST"),
            os.Getenv("DB_QUIZ_PORT"),
            os.Getenv("DB_QUIZ_NAME"),
        )
        return gorm.Open(mysql.Open(qsn), &gorm.Config{})
}