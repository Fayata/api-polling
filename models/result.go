package models

type Result struct{
	ID int `json:"id"`
	Vote int `json:"vote"`
	Participant int `json:"particpiant"`

}
type PollingResult struct {
    Title             string  `json:"title"`
    Item1Count        int     `json:"item1_count"`
    Item2Count        int     `json:"item2_count"`
    Item1Percentage   float64 `json:"item1_percentage"`
    Item2Percentage   float64 `json:"item2_percentage"`
    TotalParticipants int     `json:"total_participants"`
}
var Res = map[int]*Result{}