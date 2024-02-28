package dto

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/util"
	"github.com/ziadrahmatullah/minimarket-app/valueobject"
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

type BestCategoriesRes struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type OrderHistoryParam struct {
	SortBy *string `form:"sort_by" binding:"omitempty,oneof=order_date price"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (op *OrderHistoryParam) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()

	if op.Page != nil {
		query.WithPage(*op.Page)
	}
	if op.Limit != nil {
		query.WithLimit(*op.Limit)
	}

	if op.Order != nil {
		query.WithOrder(valueobject.Order(*op.Order))
	}

	if op.SortBy != nil {
		query.WithSortBy(*op.SortBy)
	}

	return query, nil
}

type OrderItemResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	SubTotal string `json:"sub_total"`
}

type OrderHistoryResponse struct {
	Id            string               `json:"id"`
	OrderItem     []*OrderItemResponse `json:"order_items"`
	OrderDate     string               `json:"order_date"`
	TotalPayment  string               `json:"total_payment"`
	Payment       string               `json:"payment"`
	PaymentReturn string               `json:"payment_return"`
	PaymentMethod string               `json:"payment_method"`
}

type ReportDailyQueryParamReq struct {
	Date   *string `form:"date"`
	SortBy *string `form:"sort_by" binding:"omitempty,oneof=date"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (p *ReportDailyQueryParamReq) ToQuery() *valueobject.Query {
	query := valueobject.NewQuery()

	if p.Page != nil {
		query.WithPage(*p.Page)
	}
	if p.Limit != nil {
		query.WithLimit(*p.Limit)
	}

	if p.Order != nil {
		query.WithOrder(valueobject.Order(*p.Order))
	}

	if p.SortBy != nil {
		query.WithSortBy(*p.SortBy)
	} else {
		query.WithSortBy("id")
	}

	if p.Date != nil {
		date, _ := util.ParseDate(*p.Date)
		startOfDay := date.Truncate(24 * time.Hour)
		endOfDay := startOfDay.Add(24 * time.Hour)
		query.Condition("ordered_at", valueobject.GreaterThanEqual, startOfDay).
			Condition("ordered_at", valueobject.LessThan, endOfDay)
	}
	return query
}

type DailyReportReq struct{
	Date string `json:"date" binding:"required"`
}

type MonthlyReportReq struct{
	Date string `json:"date" binding:"required"`
}

