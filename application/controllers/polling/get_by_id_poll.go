package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type PollResponse struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Question    string                 `json:"question"`
	Option      map[string]interface{} `json:"option"`
	Banner      map[string]interface{} `json:"banner"`
	IsSubmitted bool                   `json:"is_submitted"`
	IsEnded     bool                   `json:"is_ended"`
	Choices     []PollChoiceResponse   `json:"choices"`
}

type PollChoiceResponse struct {
	ID         int     `json:"id"`
	Label      string  `json:"label"`
	ImageURL   string  `json:"image_url"`
	Value      float32 `json:"value"`
	IsSelected bool    `json:"is_selected"`
}

var db *gorm.DB
var dbq *gorm.DB

func init() {
	pollingDB, quizDB := database.InitDB()
	db = pollingDB.Db
	dbq = quizDB.Db
}

func ByID(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Poll_id tidak valid")
	}

	var polling models.Poll
	var userAnswer models.User_Answer

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}
	err = db.First(&polling, id).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Polling not found")
	}
	userID := e.Get("user_id").(int)

	//Fetch the poll data
	polling.ID = id
	userAnswer.ID = userID
	
	if err := db.First(&polling, polling).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Polling tidak ditemukan")
	}

	var pollChoices []models.Poll_Choices
	if err := db.Raw("SELECT * FROM poll_choices WHERE poll_id = ?", id).Scan(&pollChoices).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil pilihan polling")
	}

	var pollResults []models.Poll_Result
	if err := db.Raw("SELECT poll_id FROM poll_result WHERE poll_id = ?", id).Scan(&pollResults).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil hasil polling")
	}

	// userIDInterface := e.Get("user_id")
	var userAnswers []models.User_Answer
	if err := db.Raw("SELECT user_id, poll_id FROM user_answer WHERE user_id = ? AND poll_id = ?", userID, id).Scan(&userAnswers).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil jawaban pengguna")
	}
	// 	}
	// }

	// Check if poll is submitted and ended
	isSubmitted, err := models.IsSubmittedPoll(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking submission status")
	}

	isEnded, err := models.IsEndedPoll(polling.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking poll end status")
	}

	// Format choices for response
	formattedChoices := make([]PollChoiceResponse, len(pollChoices))
	for i, choice := range pollChoices {
		isSelected := false
		for _, ua := range userAnswers {
			if ua.Choice_id == choice.ID {
				isSelected = true
				break
			}
		}
		votePercentage, err := choice.GetVotePercentage(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil persentase suara")
		}

		formattedChoices[i] = PollChoiceResponse{
			ID:         choice.ID,
			Label:      choice.Choice_text,
			ImageURL:   choice.Choice_image,
			Value:      votePercentage,
			IsSelected: isSelected,
		}
	}
	err = polling.GetByID(polling.ID)

	//Get choice type
	choiceType := models.GetChoiceType(polling.Options_type)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mengambil tipe pilihan")
	}
	bannerType := models.GetBannerType(polling.Question_image)
	//Get meta data
	metaData := database.Meta()

	message := "Success"
	code := 0
	if err != nil {
		message = err.Error()
		code = 1
	}

	response := PollResponse{
		ID:       polling.ID,
		Title:    polling.Title,
		Question: polling.Question_text,
		Option: map[string]interface{}{
			"type": choiceType,
			"data": formattedChoices,
		},
		Banner: map[string]interface{}{
			"type": bannerType,
			"url":  polling.Question_image,
		},
		IsSubmitted: isSubmitted,
		IsEnded:     isEnded,
		Choices:     formattedChoices,
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"data": response,
		"meta": map[string]interface{}{
			"image_path": metaData.ImagePath,
			"video_path": metaData.VideoPath,
		},
		"status": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	})
}
