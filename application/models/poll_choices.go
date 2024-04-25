package models

type PollChoices struct {
	ID      int    `json:"id"`
	Poll_id int    `json:"poll_id"`
	Option  string `json:"option"`
	Title [] PollingTitle `json:"titel"`
}

type PollingTitle struct {
	ID      int          `json:"id"`
	Title   string       `json:"title"`
}