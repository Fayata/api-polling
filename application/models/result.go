package models

type Result struct {
	ID          int `json:"id"`
	Vote        int `json:"vote"`
	Participant int `json:"particpiant"`
	User_id     int `json:"user_id"`
	Poll_id     int `json:"poll_id"`
}
type PollingResult struct {
	Title             string  `json:"title"`
	Item1Percentage   float64 `json:"item1_percentage"`
	Item2Percentage   float64 `json:"item2_percentage"`
	Item3Percentage   float64 `json:"item3_percentage"`
	Item4Percentage   float64 `json:"item4_percentage"`
	Item5Percentage   float64 `json:"item5_percentage"`
	TotalParticipants int     `json:"total_participants"`
}

var Res = map[int]*Result{}
