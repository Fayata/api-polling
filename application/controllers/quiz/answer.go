package controllers

import (
	"api-polling/application/models"

	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type UserAnswersQuiz struct {
	UserID     int `gorm:"column:user_id;NOT NULL"`
	ChoiceID   int `gorm:"column:choice_id;NOT NULL"`
	QuestionID int `gorm:"column:question_id;NOT NULL"`
	QuizID     int `gorm:"column:quiz_id;NOT NULL"`
}

type UserAnswersQuizDTO struct {
	ChoiceID   int `json:"choice_id"`
	QuestionID int `json:"question_id"`
}
type ResponseQuiz struct {
	Status struct {
		Code    int    ` json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"status"`

	Data struct{} `json:"data"`
	Meta struct{} `json:"meta"`
}

func Answer(e echo.Context) error {
	var (
		response  ResponseQuiz
		answerArr []UserAnswersQuizDTO
	)

	userID := e.Get("user_id").(int)
	userChoice := models.UserAnswer{
		UserID: userID,
	}
	if err := e.Bind(&answerArr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	quizIDStr := e.Param("id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quiz_id")
	}
	userChoice.QuizID = int(quizID)

	// Panggil fungsi CheckQuizStatus model
	status, err := models.CheckQuizStatus(userID, quizID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid cek quiz status")
	}

	if status {
		response.Status.Message = "User has alredy Quiz"
		response.Status.Type = "error"
		response.Status.Code = 1
		return e.JSON(http.StatusBadRequest, response)
	}
	// Panggil fungsi AddPoll dari model
	// err = userChoice.AddQuiz()
	for _, answer := range answerArr {
		err := models.AddQuiz(userID, quizID, answer.QuestionID, uint(answer.ChoiceID))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	response.Status.Message = "Submit Success"
	response.Status.Type = "Success"
	response.Status.Code = 0
	return e.JSON(http.StatusOK, response)
}
