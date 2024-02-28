package repository

import (
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository[entity.Product]
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
