package models

type Polling struct{
	Poll_id int
	Title string
	Item1 string
	Item2 string
}
var Poll = map[int]*Polling{}
