package repository

import (
	"context"
	"strings"
	"time"

	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/logger"
	"github.com/ziadrahmatullah/minimarket-app/util"
	"github.com/ziadrahmatullah/minimarket-app/valueobject"
	"gorm.io/gorm"
)

type OrderRepository interface {
	BaseRepository[entity.Order]
	DailyOrderReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	FindAllOrders(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type orderRepository struct {
	*baseRepository[entity.Order]
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db:             db,
		baseRepository: &baseRepository[entity.Order]{db: db},
	}
}

func (r *orderRepository) DailyOrderReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "date":
			query.WithSortBy("\"orders\".ordered_at")
		case "id":
			query.WithSortBy("\"orders\".id ")
		}
		dateN := query.GetConditionValue("date")
		logger.Log.Info(dateN)
		if dateN != nil {
			date, _ := util.ParseDate(dateN.(string))
			startOfDay := date.Truncate(24 * time.Hour)
			endOfDay := startOfDay.Add(24 * time.Hour)
			db.Where("\"orders\".ordered_at >= ? AND \"orders\".ordered_at < = ?", startOfDay, endOfDay)
		}
		return db
	})
}

func (r *orderRepository) FindAllOrders(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "price":
			query.WithSortBy("\"total_payment\"")
		case "order_date":
			query.WithSortBy("\"ordered_at\"")
		}
		db.Joins("LEFT JOIN order_items ON orders.id = order_items.order_id").
			Joins("LEFT JOIN products ON products.id = order_items.product_id")

		db.Group("orders.id")
		db.Preload("OrderItems.Product")
		return db
	})
}
