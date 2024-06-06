package polling

import "time"

type Poll struct {
	ID             int    `gorm:"column:id"`
	Title          string `gorm:"column:title"`
	Question_text  string `gorm:"column:question_text"`
	Question_image string `gorm:"column:question_image"`
	ImageURL       string `gorm:"column:image_url"`
	Options_type   string `gorm:"column:options_type"`
	// Banner_type    string     `gorm:"column:banner_type"`
	Start_date time.Time `gorm:"column:start_date"`
	End_date   time.Time `gorm:"column:end_date"`
}

func (m *Poll) TableName() string {
	return "poll"
}
