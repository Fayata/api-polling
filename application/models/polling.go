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

//Menampilkan polling berdasarkan id
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

    // Get choices for this polling
    rows, err := db.Query("SELECT id, option FROM poll_choices WHERE poll_id = ?", p.ID)
    if err != nil {
        log.Println("Failed to execute query for choices:", err)
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var choice PollChoice
        if err := rows.Scan(&choice.ID, &choice.Option); err != nil {
            log.Println("Failed to read row from query result:", err)
            continue
        }
        p.Choices = append(p.Choices, choice)
    }

    if err := rows.Err(); err != nil {
        log.Println("Failed to read all rows from query result:", err)
        return err
    }

    return nil
}

//Menampilkan semua Polling
func (p *Polling) GetAll() ([]Polling, error) {
    var polls []Polling

    db, err := database.Conn()
    if err != nil {
        log.Println("Failed to connect to database:", err)
        return polls, err
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, title FROM polling")
    if err != nil {
        log.Println("Failed to execute query:", err)
        return polls, err
    }
    defer rows.Close()

    for rows.Next() {
        var poll Polling
        if err := rows.Scan(&poll.ID, &poll.Title); err != nil {
            log.Println("Failed to read row from query result:", err)
            continue
        }
        polls = append(polls, poll)
    }

    if err := rows.Err(); err != nil {
        log.Println("Failed to read all rows from query result:", err)
        return polls, err
    }

    return polls, nil
}