package models

type Polling struct{
	Poll_id int `json:"poll_id"`
	Title string `json:"title"`
	Item1 string `json:"item1"`
	Item2 string `json:"item2"`
}

// type PollingResponse struct {
// 	ID    int             `json:"id"`
// 	Title string          `json:"title"`
// 	Items []*PollingItem `json:"item"`
// }

// type PollingItem struct {
// 	Value string `json:"value"`
// 	Label string `json:"label"`
// }