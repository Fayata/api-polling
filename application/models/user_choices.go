package models

import (
	"api-polling/system/database"
	"errors"
	"log"
	"sync"

	"github.com/labstack/echo"
)

type UserChoice struct {
	ID        int `json:"id"`
	Choice_ID int `json:"choice_id"`
	User_id   int `json:"user_id"`
	// Poll_id   int `json:"poll_id"`
}

func (uc *UserChoice) AddPoll(e echo.Context) error {
	db, err := database.Conn()
	if err != nil {
		log.Println("Gagal terhubung ke database:", err)
		return err
	}
	defer db.Close()

	// Mendapatkan user_id dari JWT token
	userID := e.Get("user_id").(int)
	uc.User_id = userID

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Cek apakah user_id sudah pernah melakukan polling
		checkQuery := `SELECT COUNT(*) FROM user_choice WHERE user_id = ?`
		var count int
		err = db.QueryRow(checkQuery, uc.User_id).Scan(&count)
		if err != nil {
			log.Println("Gagal melakukan pengecekan polling:", err)
			errChan <- err
			return
		}

		// Jika user_id sudah melakukan polling, kembalikan error
		if count > 0 {
			errChan <- errors.New("User sudah melakukan polling")
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Jika user_id belum melakukan polling, lanjutkan dengan INSERT
		insertQuery := `
			INSERT INTO user_choice (choice_id, user_id)
			VALUES (?, ?)
		`
		_, err = db.Exec(insertQuery, uc.Choice_ID, uc.User_id)
		if err != nil {
			log.Println("Gagal menambahkan polling:", err)
			errChan <- err
			return
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
