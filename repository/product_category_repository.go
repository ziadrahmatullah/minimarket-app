package repository

import (
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"gorm.io/gorm"
)

type ProductCategoryRepository interface {
	BaseRepository[entity.ProductCategory]
}

type productCategoryRepository struct {
	*baseRepository[entity.ProductCategory]
	db *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) ProductCategoryRepository {
	return &productCategoryRepository{
		db:             db,
		baseRepository: &baseRepository[entity.ProductCategory]{db: db},
	}
}
