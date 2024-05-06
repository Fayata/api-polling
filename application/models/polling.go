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
func (p *Polling) GetByID(id int) error {
    db, err := database.Conn()
    if err != nil {
        log.Println("Failed to connect to database:", err)
        return err
    }
    defer db.Close()

    query := `
        SELECT p.id, p.title, pc.id, pc.option 
        FROM polling p
        LEFT JOIN poll_choices pc ON p.id = pc.poll_id
        WHERE p.id = ?
    `
    rows, err := db.Query(query, id)
    if err != nil {
        log.Println("Failed to execute query:", err)
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var choice PollChoice
        if err := rows.Scan(&p.ID, &p.Title, &choice.ID, &choice.Option); err != nil {
            log.Println("Failed to read row from query result:", err)
            continue
        }
        p.Choices = append(p.Choices, choice)
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


