package repository

import (
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	BaseRepository[entity.OrderItem]
}

type orderItemRepository struct {
	*baseRepository[entity.OrderItem]
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{
		db:             db,
		baseRepository: &baseRepository[entity.OrderItem]{db: db},
	}
}
