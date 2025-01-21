package models

import (
	"time"
)

type TransferTransaction struct {
	ID            uint      `gorm:"primaryKey"`
	TransactionID uint      `gorm:"not null;constraint:OnDelete:CASCADE"`
	TargetUserID  uint      `form:"sendto" gorm:"not null;constraint:OnDelete:CASCADE"`
	IsDeleted     bool      `gorm:"default:false"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
