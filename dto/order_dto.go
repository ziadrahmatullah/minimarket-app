package dto

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/entity"
)

type AddOrderReq struct {
	ProductCodes  []string `json:"product_codes" binding:"required"`
	ProductQty    []int    `json:"product_qty" binding:"required,dive,number"`
	PaymentMethod string   `json:"payment_method" binding:"required"`
	Payment       string   `json:"payment" binding:"required"`
}

type AddOrderRes struct {
	TotalPayment  string `json:"total_payment"`
	PaymentReturn string `json:"payment_return"`
}

func (r *AddOrderReq) ToOrder() (*entity.Order, error) {
	Payment, _ := decimal.NewFromString(r.Payment)
	if Payment.LessThan(decimal.Zero) {
		return nil, apperror.NewInvalidPathQueryParamError(errors.New("payment should be greater than zero"))
	}
	return &entity.Order{
		OrderedAt:     time.Now(),
		Payment:       Payment,
		PaymentMethod: r.PaymentMethod,
	}, nil
}
