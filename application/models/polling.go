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
	PollID int    `json:"poll_id"`
}

///////////////////CMS////////////////////

func (p *Polling) Create() error {
    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return err
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        log.Println("Gagal memulai transaksi:", err)
        return err
    }

    query := "INSERT INTO polling (title) VALUES (?)"
    result, err := tx.Exec(query, p.Title)
    if err != nil {
        tx.Rollback()
        log.Println("Gagal membuat polling baru (tabel polling):", err)
        return err
    }

    lastInsertId, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        log.Println("Gagal mendapatkan last insert ID:", err)
        return err
    }

    query2 := "INSERT INTO poll_choices (option, poll_id) VALUES (?, ?)"
    for _, choice := range p.Choices {
        _, err = tx.Exec(query2, choice.Option, lastInsertId)
        if err != nil {
            tx.Rollback()
            log.Println("Gagal menambahkan pilihan polling:", err)
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        log.Println("Gagal melakukan commit:", err)
        return err
    }

    return nil
}

func (p *Polling) Update(id int) error {
    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return err
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        log.Println("Gagal memulai transaksi", err)
        return err
    }

    result, err := tx.Exec("UPDATE polling SET title = ? WHERE id = ?", p.Title, id)
    if err != nil {
        tx.Rollback()
        log.Println("Gagal mengupdate data polling:", err)
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        tx.Rollback()
        log.Println("Gagal mendapatkan jumlah baris yang terdampak:", err)
        return err
    }

    if rowsAffected == 0 {
        tx.Rollback()
        log.Println("Data polling tidak ditemukan:", id)
        return err
    }

    for _, choice := range p.Choices {
        _, err = tx.Exec("UPDATE poll_choices SET option = ? WHERE id = ?", choice.Option, choice.ID)
        if err != nil {
            tx.Rollback()
            log.Println("Gagal mengupdate pilihan polling:", err)
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        log.Println("Gagal melakukan commit:", err)
        return err
    }

    return nil
}

func (p *Polling) Delete(id int) error {
    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return err
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        log.Println("Gagal memulai transaksi:", err)
        return err
    }

    result, err := tx.Exec("DELETE FROM poll_choices WHERE poll_id=?", id)
    if err != nil {
        tx.Rollback()
        log.Println("Gagal menghapus hasil polling:", err)
        return err
    }

    affectedRows, err := result.RowsAffected()
    if err != nil {
        tx.Rollback()
        log.Println("Gagal mendapatkan jumlah baris yang terhapus:", err)
        return err
    }

    if affectedRows == 0 {
        tx.Rollback()
        log.Println("Tidak ada hasil polling yang ditemukan untuk ID:", id)
        return err
    }

    _, err = tx.Exec("DELETE FROM polling WHERE id=?", id)
    if err != nil {
        tx.Rollback()
        log.Println("Gagal menghapus data polling:", err)
        return err
    }

    affectedRows, err = result.RowsAffected()
    if err != nil {
        tx.Rollback()
        log.Println("Gagal mendapatkan jumlah baris yang terhapus:", err)
        return err
    }

    if affectedRows == 0 {
        tx.Rollback()
        log.Println("Tidak ada data polling yang ditemukan untuk ID:", id)
        return err
    }

    if err := tx.Commit(); err != nil {
        log.Println("Gagal melakukan commit:", err)
        return err
    }

    return nil
}

////////////////////USERS///////////////////

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
