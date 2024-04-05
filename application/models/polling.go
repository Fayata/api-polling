package models

type Polling struct {
	Poll_id int    `json:"poll_id"`
	Title   string `json:"title"`
	Item1   string `json:"item1"`
	Item2   string `json:"item2"`
	Item3   string `json:"item3"`
	Item4   string `json:"item4"`
	Item5   string `json:"item5"`
}

type PollingResponse struct {
	ID    int             `json:"id"`
	Title string          `json:"title"`
	Items []*PollingItem `json:"items"`
}

type PollingItem struct {
	Value string `json:"value"`
}

