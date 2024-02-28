package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	Id                uint   `gorm:"primaryKey;autoIncrement"`
	Name              string `gorm:"not null"`
	ProductCategoryId uint   `gorm:"not null"`
	ProductCategory   ProductCategory
	ProductCode       string          `gorm:"not null"`
	Stock             int             `gorm:"not null"`
	Price             decimal.Decimal `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}
