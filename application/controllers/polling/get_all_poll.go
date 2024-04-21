package controllers

import (
    "api-polling/application/models"
    "api-polling/system/database"
    "log"
    "net/http"

    "github.com/labstack/echo"
)

type responseData struct {
    Data []models.Polling `json:"data"`
}

func AllList(e echo.Context) error {
    var pollList []models.Polling

    db, err := database.Conn()
    if err != nil {
        log.Println("Failed to connect to database:", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
    }
    defer db.Close()

    rows, err := db.Query("SELECT poll_id, title FROM polling")
    if err != nil {
        log.Println("Failed to execute query:", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to execute query")
    }
    defer rows.Close()

    // Iterasi melalui hasil query dan memasukkan data polling ke dalam slice pollList
    for rows.Next() {
        var polling models.Polling
        if err := rows.Scan(&polling.ID, &polling.Title); err != nil {
            log.Println("Failed to read row:", err)
            continue
        }
        pollList = append(pollList, polling)
    }

    // Jika ada kesalahan saat iterasi
    if err := rows.Err(); err != nil {
        log.Println("Failed to read all rows:", err)
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read all rows")
    }

    // Mengatur hasil sesuai dengan format yang diinginkan
    responseData := responseData{
        Data: pollList,
    }

    return e.JSON(http.StatusOK, responseData)
}
