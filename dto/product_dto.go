package dto

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/entity"
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
