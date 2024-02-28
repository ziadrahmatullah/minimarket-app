package repository

import (
	"context"
	"strings"

	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/valueobject"
	"gorm.io/gorm"
)

type ProductRepository interface {
	BaseRepository[entity.Product]
	FindAllProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
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

func (r *productRepository) FindAllProducts(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "price":
			query.WithSortBy("price")
		}

		category := query.GetConditionValue("category")
		name := query.GetConditionValue("name")
		db.Preload("ProductCategory")

		if category != nil {
			db.Where("product_category_id", category)
		}

		if name != nil {
			db.Where("products.name ILIKE ?", name)
		}

		return db
	})
}
