package database

import (
	"database/sql"
	"fmt"
	"log"
	"api-polling/application/config/app"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	initDatabase()
}

func Conn() (*sql.DB, error) {
	dbConfig := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", app.Load.Database.User, app.Load.Database.Password, app.Load.Database.Host, app.Load.Database.Port, app.Load.Database.Name)

	db, err := sql.Open("mysql", dbConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDatabase() {
	dbConfig := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", app.Load.Database.User, app.Load.Database.Password, app.Load.Database.Host, app.Load.Database.Port, app.Load.Database.Name)
	db, err := sql.Open("mysql", dbConfig)
	if err != nil {
		log.Fatalln("Error init to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("Error ping to database:", err)
	}
}
