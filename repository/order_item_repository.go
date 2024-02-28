package repository

import (
	"context"

	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	BaseRepository[entity.OrderItem]
	BulkCreate(ctx context.Context, orderItems []*entity.OrderItem) error
	GetMostOrderedCategories(ctx context.Context) ([]entity.ProductCategory, error)
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

func (r *orderItemRepository) GetMostOrderedCategories(ctx context.Context) ([]entity.ProductCategory, error) {
	var categoryCounts []entity.ProductCategory

	err := r.conn(ctx).Model(&entity.OrderItem{}).
		Select("products.product_category_id, product_categories.name AS name, COUNT(*) as order_count").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("JOIN product_categories ON products.product_category_id = product_categories.id").
		Group("products.product_category_id, product_categories.name").
		Order("order_count desc").
		Scan(&categoryCounts).
		Error
	if err != nil {
		return nil, err
	}

	return categoryCounts, nil
}
