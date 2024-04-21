package models

type UserChoice struct {
	ID int    `json:"id"`
	Choice_ID int `json:"choice_id"`
	User_id int `json:"user_id"`
}

type PollingResult struct {
	ChoiceID  int `json:"choice_id"`
	VoteCount int `json:"vote_count"`
}