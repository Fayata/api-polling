package models

type Polling struct{
	Poll_id int `json:"poll_id"`
	Title string `json:"title"`
	Item1 string `json:"item1"`
	Item2 string `json:"item2"`
}
var Poll = map[int]*Polling{

}
 