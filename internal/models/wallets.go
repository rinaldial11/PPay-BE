package models

import (
	"time"
)

type Wallet struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id" db:"user_id" gorm:"not null;unique;constraint:OnDelete:CASCADE"`
	Balance   float64   `gorm:"type:decimal;default:0.00"`
	IsDeleted bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
