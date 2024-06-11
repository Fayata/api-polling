package controllers

import (
	"api-polling/application/models"
	"api-polling/system/database"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type QuestionData struct {
	ID            int    `json:"id"`
	QuizID        int    `json:"quiz_id"`
	Number        int    `json:"number"`
	QuestionText  string `json:"question_text"`
	QuestionImage string `json:"question_image"`
	Choices []models.QuizQuestionChoice `json:"choices"`
}

func GetQuizByID(e echo.Context) error {
	id, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
	}

	var quiz models.Quiz
	var question models.QuizQuestion

	// Get quiz by ID
	db, err := database.InitDB().DbQuiz()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Quiz tidak ditemukan")
	}
	err = db.First(&quiz, id).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Quiz tidak ditemukan")
	}
	err = db.First(&question, id).Error
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Option tidak ditemukan")
	}
	userID := e.Get("user_id").(int)

	// Check if quiz is submitted and ended
	isSubmitted, isEnded, err := quiz.CheckQuizStatus(e, uint(userID))
	if err != nil {
		// Tangani kesalahan jika ada
		log.Println("Gagal memeriksa status quiz:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal memeriksa status quiz")
	}

	currentQuestion, err := quiz.GetQuizPosition()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mendapatkan posisi quiz")
	}
	//Mengambil data question

	questionDataArr, err := models.GetQuestionByQuizId(id)
	if err != nil {
		log.Println("Gagal mendapatkan question")
		return echo.NewHTTPError(http.StatusInternalServerError, "Gagal mendapatkan posisi question")
	}

	var questionResults []QuestionData
	for _, value := range questionDataArr {
		var questionResult QuestionData
		questionResult.ID = value.ID
		questionResult.Number = value.Number
		questionResult.QuestionImage = value.QuestionImage
		questionResult.QuestionText = value.QuestionText
		questionResult.QuizID = value.QuizID

		// get question choices
		choiceDataArr, err := models.GetChoiceByQuestionId(value.ID)
		if err != nil {
			log.Println("failed to get question choices")
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get question choices")
		}
		for index, valueChoice := range choiceDataArr {
			choiceDataArr[index].QuestionType = models.GetQuestionType(valueChoice.ChoiceImage)
		}
		questionResult.Choices = choiceDataArr
		questionResults = append(questionResults, questionResult)
	}

	// Format response sesuai dengan struktur yang diinginkan
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"id":       quiz.ID,
			"title":    quiz.Name,
			"question": questionResults,

			"banner": map[string]interface{}{
				"type": quiz.IsActive,
				"url":  question.QuestionImage,
			},
			"is_submitted": isSubmitted,
			"is_ended":     isEnded,
		},
		"meta": map[string]interface{}{
			"questions": map[string]interface{}{
				"total":   quiz.TotalQuestion,
				"current": currentQuestion,
			},
		},
		"status": map[string]interface{}{
			"code":    0,
			"message": "Success",
		},
	}
	// if err != nil {
	// 	response["status"].(map[string]interface{})["code"] = 1
	// 	response["status"].(map[string]interface{})["message"] = err.Error()
	// 	if err.Error() == "user has already polled" {
	// 		return e.JSON(http.StatusBadRequest, response)
	// 	} else {
	// 		return e.JSON(http.StatusInternalServerError, response)
	// 	}
	// }

	return e.JSON(http.StatusOK, response)
}
