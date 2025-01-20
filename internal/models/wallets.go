package models

import (
	"time"
)

type Wallet struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;unique;constraint:OnDelete:CASCADE"`
	Balance   float64   `gorm:"type:decimal(10,2);default:0.00;check:balance >= 0"`
	IsDeleted bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
