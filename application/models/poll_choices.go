package models

import (
	"api-polling/system/database"
	"log"
)

type PollChoices struct {
	ID     int    `json:"id"`
	PollID int    `json:"poll_id"`
	Option string `json:"option"`
}

type PollingTitle struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type UserChoices struct {
	ID        int `json:"id"`
	Choice_ID int `json:"choice_id"`
	User_id   int `json:"user_id"`
}

type PollingResult struct {
	ChoiceID   int     `json:"choice_id"`
	VoteCount  int     `json:"vote_count"`
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
        SELECT id, title
        FROM polling
        WHERE id = ?
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

	totalVotes, err := getTotalVotes(pollID)
	if err != nil {
		log.Println("Gagal mengambil total suara:", err)
		return nil, err
	}

	for rows.Next() {
		var choiceID int
		var option string
		var voteCount int

		err = rows.Scan(&choiceID, &option, &voteCount)
		if err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}

		percentage := float64(voteCount) / float64(totalVotes) * 100

		result := map[string]interface{}{
			"poll_choice": map[string]interface{}{
				"id":     choiceID,
				"option": option,
				"title":  pollingTitle,
			},
			"percentage": percentage,
		}

		pollingResults = append(pollingResults, result)
	}

	return pollingResults, nil
}

// Fungsi untuk menghitung total suara
func getTotalVotes(pollID int) (int, error) {
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return 0, err
	}
	defer db.Close()

	var totalVotes int
	err = db.QueryRow("SELECT COUNT(*) FROM user_choice WHERE choice_id IN (SELECT id FROM poll_choices WHERE poll_id = ?)", pollID).Scan(&totalVotes)
	if err != nil {
		log.Println("Gagal menjalankan query:", err)
		return 0, err
	}
	return totalVotes, nil
}
