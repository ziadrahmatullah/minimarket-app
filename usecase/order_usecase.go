package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/repository"
	"github.com/ziadrahmatullah/minimarket-app/transactor"
	"github.com/ziadrahmatullah/minimarket-app/util"
	"github.com/ziadrahmatullah/minimarket-app/valueobject"
)

type OrderUsecase interface {
	AddOrder(ctx context.Context, order *entity.Order, productCodes []string, productQty []int) (*entity.Order, error)
	GetMostOrderedCategories(ctx context.Context) ([]entity.ProductCategory, error)
	DailyOrderReport(ctx context.Context, date time.Time) ([]*entity.Order, error)
	ListAllOrders(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type orderUsecase struct {
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	productRepo   repository.ProductRepository
	manager       transactor.Manager
}

func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	productRepo repository.ProductRepository,
	manager transactor.Manager,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		productRepo:   productRepo,
		manager:       manager,
	}
}

func (u *orderUsecase) AddOrder(ctx context.Context, order *entity.Order, productCodes []string, productQty []int) (*entity.Order, error) {
	if len(productCodes) != len(productQty) {
		return nil, apperror.NewResourceStateError("product and qty not match")
	}
	if !util.IsUnique(productCodes) {
		return nil, apperror.NewResourceStateError("product codes must be unique")
	}
	productQ := valueobject.NewQuery().Condition("product_code", valueobject.In, productCodes).Lock()
	fetchedProducts, err := u.productRepo.Find(ctx, productQ)
	if err != nil {
		return nil, err
	}
	if len(productCodes) != len(fetchedProducts) {
		return nil, apperror.NewResourceNotFoundError("products", "product_code", productCodes)
	}
	for i, product := range fetchedProducts {
		if product.Stock < productQty[i] {
			return nil, apperror.NewResourceStateError("product out of stock")
		}
	}
	err = u.manager.Run(ctx, func(c context.Context) error {
		newOrder, err := u.orderRepo.Create(c, order)
		if err != nil {
			return err
		}
		var orderItems []*entity.OrderItem
		var totalPayment decimal.Decimal
		for i, product := range fetchedProducts {
			orderItem := &entity.OrderItem{
				OrderId:   newOrder.Id,
				ProductId: product.Id,
				Quantity:  productQty[i],
				SubTotal:  product.Price.Mul(decimal.NewFromInt(int64(productQty[i]))),
			}
			totalPayment = totalPayment.Add(orderItem.SubTotal)
			orderItems = append(orderItems, orderItem)
			product.Stock -= productQty[i]
			_, err = u.productRepo.Update(c, product)
			if err != nil {
				return err
			}
		}
		newOrder.TotalPayment = totalPayment
		if order.Payment.LessThan(totalPayment) {
			return apperror.NewResourceStateError("your money is not enough")
		}
		if decimal.NewFromInt(100000).LessThan(totalPayment) {
			minumanId := 2
			productBonusQ := valueobject.NewQuery().
				Condition("product_category_id", valueobject.Equal, minumanId).
				Condition("stock", valueobject.GreaterThan, 0)
			fetchedProductBonus, err := u.productRepo.Find(c, productBonusQ)
			if err != nil {
				return err
			}
			productBonus := getRandomValueFromArray(fetchedProductBonus)
			productBonus.Stock--
			_, err = u.productRepo.Update(c, productBonus)
			if err != nil {
				return err
			}
			orderItem := &entity.OrderItem{
				OrderId:   newOrder.Id,
				ProductId: productBonus.Id,
				Quantity:  1,
				SubTotal:  decimal.Zero,
			}
			orderItems = append(orderItems, orderItem)
		}
		newOrder.PaymentReturn = order.Payment.Sub(totalPayment)
		err = u.orderItemRepo.BulkCreate(c, orderItems)
		if err != nil {
			return err
		}
		order, err = u.orderRepo.Update(c, newOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return order, err
}

func (u *orderUsecase) GetMostOrderedCategories(ctx context.Context) ([]entity.ProductCategory, error) {
	return u.orderItemRepo.GetMostOrderedCategories(ctx)
}

// func (u *orderUsecase) DailyOrderReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
// 	return u.orderRepo.DailyOrderReport(ctx, query)
// }

func (u *orderUsecase) DailyOrderReport(ctx context.Context, date time.Time) ([]*entity.Order, error) {
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	orderQuery := valueobject.NewQuery().
		Condition("ordered_at", valueobject.GreaterThanEqual, startOfDay).
		Condition("ordered_at", valueobject.LessThan, endOfDay)
	fethcedOrder, err := u.orderRepo.Find(ctx, orderQuery)
	if err != nil {
		return nil, err
	}
	return fethcedOrder, nil
}

func (u *orderUsecase) ListAllOrders(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return u.orderRepo.FindAllOrders(ctx, query)
}

func getRandomValueFromArray(arr []*entity.Product) *entity.Product {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(arr))
	return arr[randomIndex]
}
