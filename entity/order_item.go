package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderItem struct {
	Id        uint            `gorm:"primaryKey;autoIncrement"`
	OrderId   uint            `gorm:"not null"`
	Order     Order           `gorm:"foreignKey:OrderId;references:Id"`
	ProductId uint            `gorm:"not null"`
	Product   Product         `gorm:"foreignKey:ProductId;references:Id"`
	Quantity  int             `gorm:"not null"`
	SubTotal  decimal.Decimal `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
