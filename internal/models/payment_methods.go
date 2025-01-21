package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Money float64

func (m *Money) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid type for Money")
	}

	// Remove the dollar sign and parse the number
	strValue = strings.Replace(strValue, "$", "", -1)
	parsedValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return err
	}

	*m = Money(parsedValue)
	return nil
}

type PaymentMethod struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Tax       Money   `gorm:"type:money;default:0.00"`
	IsDeleted bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	TopupTransactions []TopupTransaction `gorm:"constraint:OnDelete:CASCADE"`
}
