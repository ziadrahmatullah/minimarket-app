package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order struct {
	Id            uint            `gorm:"primaryKey;autoIncrement"`
	OrderedAt     time.Time       `gorm:"not null"`
	OrderStatusId uint            `gorm:"not null"`
	ItemOrderQty  int             `gorm:"not null"`
	TotalPayment  decimal.Decimal `gorm:"not null;type:numeric"`
	Payment       decimal.Decimal `gorm:"not null;type:numeric"`
	OrderItems    []OrderItem     `gorm:"foreignKey:OrderId"`
	PaymentMethod string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
