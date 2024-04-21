package models

type Polling struct {
	ID      string       `json:"id"`
	Title   string       `json:"title"`
	Choices []PollChoice `json:"choices"`
}

type PollChoice struct {
	ID     string `json:"id"`
	Option string `json:"option"`
}
