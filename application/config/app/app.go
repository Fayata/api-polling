package app

import (
	"database/sql"
	"fmt"
	"log"
)

var Load config

type config struct {
	WebServer webServer
	Database  database
}

type webServer struct {
	Host string
	Port string
}

type database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func init() {
	Load = config{
		WebServer: webServer{
			Host: "0.0.0.0",
			Port: "9000",
		},
		Database: database{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "",
			Name:     "db_api",
		},
	}

	// === Bootstraping initiation
	initDatabase()
}

func initDatabase() {
	dbConfig := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", Load.Database.User, Load.Database.Password, Load.Database.Host, Load.Database.Port, Load.Database.Name)
	db, err := sql.Open("mysql", dbConfig)
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}
}
