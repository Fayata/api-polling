package models

type PollChoices struct{
	ID int    `json:"id"`
	Poll_id int `json:"poll_id"`
	Option string `json:"option"`
}


