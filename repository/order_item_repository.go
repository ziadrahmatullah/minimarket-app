package repository

import (
	"context"

	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	BaseRepository[entity.OrderItem]
	BulkCreate(ctx context.Context, orderItems []*entity.OrderItem) error
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

func (r *orderItemRepository) BulkCreate(ctx context.Context, orderItems []*entity.OrderItem) error {
	err := r.conn(ctx).Model(&entity.OrderItem{}).Create(orderItems).Error
	if err != nil {
		return err
	}
	return nil
}
