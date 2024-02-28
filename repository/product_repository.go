package repository

import (
	"context"

	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository[entity.Product]
	FindUnique(ctx context.Context, productCodes []string) ([]entity.Product, error)
}

type productRepository struct {
	*baseRepository[entity.Product]
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Product]{db: db},
	}
}

func (r *productRepository) FindUnique(ctx context.Context, productCodes []string) ([]entity.Product, error) {
	var products []entity.Product
	err := r.conn(ctx).Model(&entity.Product{}).Where("product_code IN ?", productCodes).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
