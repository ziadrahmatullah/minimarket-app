package dto

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/valueobject"
)

type AddProductReq struct {
	ProductCode       string `json:"product_code" binding:"required"`
	Name              string `json:"name" binding:"required"`
	ProductCategoryId uint   `json:"product_category_id" binding:"required"`
	Stock             int    `json:"stock" binding:"required"`
	Price             string `json:"price" binding:"required,numeric"`
}

func (r *AddProductReq) ToProduct() (*entity.Product, error) {
	Price, _ := decimal.NewFromString(r.Price)
	if Price.LessThan(decimal.Zero) {
		return nil, apperror.NewInvalidPathQueryParamError(errors.New("price should be greater than zero"))
	}
	return &entity.Product{
		Name:              r.Name,
		ProductCategoryId: r.ProductCategoryId,
		ProductCode:       r.ProductCode,
		Stock:             r.Stock,
		Price:             Price,
	}, nil
}

type ListProductQueryParam struct {
	Name     *string `form:"name"`
	Category *int    `form:"category" binding:"omitempty,numeric,min=1"`
	SortBy   *string `form:"sort_by" binding:"omitempty,oneof=name price"`
	Order    *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit    *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page     *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (qp *ListProductQueryParam) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()

	if qp.Page != nil {
		query.WithPage(*qp.Page)
	}
	if qp.Limit != nil {
		query.WithLimit(*qp.Limit)
	}

	if qp.Order != nil {
		query.WithOrder(valueobject.Order(*qp.Order))
	}

	if qp.SortBy != nil {
		query.WithSortBy(*qp.SortBy)
	}

	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}

	if qp.Category != nil {
		query.Condition("category", valueobject.Equal, *qp.Category)
	}

	return query, nil
}

type ProductRes struct {
	Id                uint            `json:"id"`
	Name              string          `json:"name"`
	ProductCategoryId uint            `json:"product_category_id"`
	Price             decimal.Decimal `json:"price"`
	Stock             int             `json:"stock"`
}

func NewFromProduct(product *entity.Product) *ProductRes {
	return &ProductRes{
		Id:                product.Id,
		Name:              product.Name,
		ProductCategoryId: product.ProductCategoryId,
		Price:             product.Price,
		Stock:             product.Stock,
	}
}
