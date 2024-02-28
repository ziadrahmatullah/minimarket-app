package repository

import (
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	BaseRepository[entity.Order]
}

type orderRepository struct {
	*baseRepository[entity.Order]
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Order]{db: db},
	}
}
