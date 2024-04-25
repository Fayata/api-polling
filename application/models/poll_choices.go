package models

import (
    "api-polling/system/database"
    "log"
)

type PollChoices struct {
    ID      int    `json:"id"`
    PollID  int    `json:"poll_id"`
    Option  string `json:"option"`
}

type PollingTitle struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
}

// Fungsi untuk mendapatkan hasil polling berdasarkan ID polling
func GetPollingResultsByID(pollID int) ([]map[string]interface{}, error) {
    var pollingResults []map[string]interface{}

    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return nil, err
    }
    defer db.Close()

    // Query untuk mendapatkan data polling berdasarkan ID
    queryPolling := `
        SELECT poll_id, title
        FROM polling
        WHERE poll_id = ?
    `
    row := db.QueryRow(queryPolling, pollID)
    var pollingTitle string
    if err := row.Scan(&pollID, &pollingTitle); err != nil {
        log.Println("Failed to execute query or no rows found:", err)
        return nil, err
    }

    // Query untuk mendapatkan pilihan polling berdasarkan ID polling
    queryChoices := `
        SELECT pc.id, pc.option, COUNT(uc.user_id) as vote_count
        FROM poll_choices pc
        LEFT JOIN user_choice uc ON pc.id = uc.choice_id
        WHERE pc.poll_id = ?
        GROUP BY pc.id, pc.option
    `
    rows, err := db.Query(queryChoices, pollID)
    if err != nil {
        log.Println("Failed to execute query for choices:", err)
        return nil, err
    }
    defer rows.Close()

    totalParticipants, err := getTotalParticipants(pollID)
    if err != nil {
        log.Println("Gagal mengambil total peserta polling:", err)
        return nil, err
    }

    for rows.Next() {
        var choice PollChoices
        var voteCount int
        if err := rows.Scan(&choice.ID, &choice.Option, &voteCount); err != nil {
            log.Println("Failed to read row from query result:", err)
            continue
        }
        percentage := float64(voteCount) / float64(totalParticipants) * 100

        result := map[string]interface{}{
            "poll_choice": map[string]interface{}{
				"title": pollingTitle, 
                "id":      choice.ID,
                "option":  choice.Option,
            },
            "percentage": percentage,
        }
        pollingResults = append(pollingResults, result)
    }

    return pollingResults, nil
}

// Fungsi untuk menghitung total peserta polling
func getTotalParticipants(pollID int) (int, error) {
    db, err := database.Conn()
    if err != nil {
        log.Println("Gagal terhubung ke database:", err)
        return 0, err
    }
    defer db.Close()

    var totalParticipants int
    err = db.QueryRow("SELECT COUNT(*) FROM poll_choices WHERE poll_id =?", pollID).Scan(&totalParticipants)
    if err != nil {
        log.Println("Gagal menjalankan query:", err)
        return 0, err
    }
    return totalParticipants, nil
}
