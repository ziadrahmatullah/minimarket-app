package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order struct {
	Id            uint            `gorm:"primaryKey;autoIncrement"`
	OrderedAt     time.Time       `gorm:"not null"`
	TotalPayment  decimal.Decimal `gorm:"type:numeric"`
	Payment       decimal.Decimal `gorm:"not null;type:numeric"`
	PaymentReturn decimal.Decimal `gorm:"type:numeric"`
	OrderItems    []OrderItem     `gorm:"foreignKey:OrderId"`
	PaymentMethod string          `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
