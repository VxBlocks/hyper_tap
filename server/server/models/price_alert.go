package models

type PriceAlertReadOrm struct {
	ID      uint   `gorm:"primarykey"`
	AlertId uint32 `gorm:"index"`
	UserId  string `gorm:"index"`
}

func (PriceAlertReadOrm) TableName() string { return "price_alert_read" }
