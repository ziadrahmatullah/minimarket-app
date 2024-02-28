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

type ReportMonthlyQueryParamReq struct {
	Date   *string `form:"date"`
	SortBy *string `form:"sort_by" binding:"omitempty,oneof=date"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

func (p *ReportMonthlyQueryParamReq) ToQuery() *valueobject.Query {
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
		year, month, _ := date.Date()
		startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0)
		query.Condition("ordered_at", valueobject.GreaterThanEqual, startOfMonth).
		Condition("ordered_at", valueobject.LessThan, endOfMonth)
	}
	return query
}
