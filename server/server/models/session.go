package models

import (
	"time"
)

type SessionORMV1 struct {
	Address   string `gorm:"address;primaryKey"`
	Session   string `gorm:"session;primaryKey"`
	CreatedAt time.Time
}

func (s SessionORMV1) TableName() string {
	return "sessions"
}
