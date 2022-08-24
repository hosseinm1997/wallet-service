package models

import (
	"time"
)

type Transaction struct {
	ID             uint `gorm:"primary_key"`
	Mobile         string
	UsageLogId     *uint
	CreditCodeText string
	Status         uint
	CreatedAt      time.Time
	Amount         *uint
}
