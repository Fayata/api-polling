package database

import (
	"api-polling/application/config/app"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Conn() (*sql.DB, error) {
	dbConfig := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", app.Load.Database.User, app.Load.Database.Password, app.Load.Database.Host, app.Load.Database.Port, app.Load.Database.Name)

	db, err := sql.Open("mysql", dbConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}
