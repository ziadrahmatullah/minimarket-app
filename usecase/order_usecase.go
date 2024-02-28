package usecase

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/repository"
	"github.com/ziadrahmatullah/minimarket-app/transactor"
	"github.com/ziadrahmatullah/minimarket-app/util"
)

type OrderUsecase interface {
	AddOrder(ctx context.Context, order *entity.Order, productCodes []string, productQty []int) (*entity.Order, error)
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
		return nil,apperror.NewResourceStateError("product codes must be unique")
	}
	fetchedProducts, err := u.productRepo.FindUnique(ctx, productCodes)
	if err != nil {
		return nil,err
	}
	if len(productCodes) != len(fetchedProducts) {
		return nil,apperror.NewResourceNotFoundError("products", "product_code", productCodes)
	}
	for i, product := range fetchedProducts {
		if product.Stock < productQty[i] {
			return nil,apperror.NewResourceStateError("product out of stock")
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
		}
		newOrder.TotalPayment = totalPayment
		if order.Payment.LessThan(totalPayment) {
			return apperror.NewResourceStateError("your money is not enough")
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
