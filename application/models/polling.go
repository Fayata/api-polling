package models

import (
    "api-polling/system/database"
    "log"
)

type Polling struct {
    ID      int          `json:"id"`
    Title   string       `json:"title"`
    Choices []PollChoice `json:"choices"`
}
type PollChoice struct {
    ID     int    `json:"id"`
    Option string `json:"option"`
}
// func (up *Polling) GetByID(id int) error {
//     db, err := database.Conn()
//     if err != nil {
//         log.Println("Gagal terhubung ke database:", err)
//         return err
//     }
//     defer db.Close()

//     query := `
//         SELECT id, title FROM polling WHERE id = ?
//     `
//     err = db.QueryRow(query, up.ID).Scan(&up.ID, &up.Title)
//     if err != nil {
//         log.Println("Gagal mengeksekusi query atau tidak ada baris yang ditemukan:", err)
//         return err
//     }

//     return nil
// }
func (p *Polling) GetByID(id int) error {
    db, err := database.Conn()
    if err != nil {
        log.Println("Failed to connect to database:", err)
        return err
    }
    defer db.Close()

    err = db.QueryRow("SELECT id, title FROM polling WHERE id = ?", id).
        Scan(&p.ID, &p.Title)
    if err != nil {
        log.Println("Failed to execute query or no rows found:", err)
        return err
    }
    return nil
}

func (up *Polling) GetAll() ([]Polling, error) {
    var polls []Polling

    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return polls, err
    }
    defer db.Close()

    query := `
        SELECT id, title FROM polling
    `
    rows, err := db.Query(query)
    if err != nil {
        log.Println("Gagal mengeksekusi query:", err)
        return polls, err
    }
    defer rows.Close()

    for rows.Next() {
        var poll Polling
        if err := rows.Scan(&poll.ID, &poll.Title); err != nil {
            log.Println("Gagal membaca baris dari hasil query:", err)
            continue
        }
        polls = append(polls, poll)
    }

    if err := rows.Err(); err != nil {
        log.Println("Gagal membaca semua baris dari hasil query:", err)
        return polls, err
    }

    return polls, nil
}


