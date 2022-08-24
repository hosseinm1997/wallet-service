package models

import (
	"time"
)

type User struct {
	Mobile       string `gorm:"primary_key"`
	Balance      uint
	CreatedAt    time.Time
	Transactions []*Transaction `gorm:"foreignKey:mobile;references:mobile"`
}
