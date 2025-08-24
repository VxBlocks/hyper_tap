package models

import "github.com/lib/pq"

type FcmToken struct {
	UserId string         `gorm:"primaryKey"`
	Tokens pq.StringArray `gorm:"type:text[]"`
}

func (f FcmToken) TableName() string {
	return "fcm_tokens"
}
