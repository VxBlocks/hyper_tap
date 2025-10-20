package models

type NewsReadOrm struct {
	ID     uint   `gorm:"primarykey"`
	NewsId string `gorm:"index"`
	UserId string `gorm:"index"`
}

func (NewsReadOrm) TableName() string { return "news_read" }
