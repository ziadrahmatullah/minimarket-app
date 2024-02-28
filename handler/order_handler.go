package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/minimarket-app/dto"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/usecase"
	"github.com/ziadrahmatullah/minimarket-app/util"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(usecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: usecase}
}

func (h *OrderHandler) AddOrder(c *gin.Context) {
	var request dto.AddOrderReq
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	order, err := request.ToOrder()
	if err != nil {
		_ = c.Error(err)
		return
	}
	order, err = h.usecase.AddOrder(c.Request.Context(), order, request.ProductCodes, request.ProductQty)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.AddOrderRes{
		TotalPayment:  order.TotalPayment.String(),
		PaymentReturn: order.PaymentReturn.String()}})
}

func (h *OrderHandler) GetMostOrderedCategories(c *gin.Context) {
	categories, err := h.usecase.GetMostOrderedCategories(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	var result []*dto.BestCategoriesRes
	for _, category := range categories {
		result = append(result, &dto.BestCategoriesRes{
			Name:  category.Name,
			Count: category.OrderCount,
		})
	}
	c.JSON(http.StatusOK, dto.Response{Data: result})
}

func (h *OrderHandler) DailyOrderReport(c *gin.Context) {
	var request dto.DailyRepotReq
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	date, err := util.ParseDate(request.Date)
	if err != nil {
		_ = c.Error(err)
		return
	}
	orders, err := h.usecase.DailyOrderReport(c.Request.Context(), date)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: orders})
}

// func (h *OrderHandler) DailyOrderReport(c *gin.Context) {
// 	var requestParam dto.ReportDailyQueryParamReq
// 	if err := c.ShouldBindQuery(&requestParam); err != nil {
// 		_ = c.Error(err)
// 		return
// 	}
// 	query := requestParam.ToQuery()
// 	pageResult, err := h.usecase.DailyOrderReport(c.Request.Context(), query)
// 	if err != nil {
// 		_ = c.Error(err)
// 		return
// 	}
// 	orders := pageResult.Data.([]*entity.Order)
// 	c.JSON(http.StatusOK, dto.Response{
// 		Data:        orders,
// 		TotalPage:   &pageResult.TotalPage,
// 		TotalItem:   &pageResult.TotalItem,
// 		CurrentPage: &pageResult.CurrentPage,
// 		CurrentItem: &pageResult.CurrentItems,
// 	})
// }

func (h *OrderHandler) OrderHistory(c *gin.Context) {
	var request dto.OrderHistoryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pagedResult, err := h.usecase.ListAllOrders(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var listOrders []*dto.OrderHistoryResponse
	data := pagedResult.Data.([]*entity.Order)
	var orders *dto.OrderHistoryResponse
	for i, item := range data {
		orders = &dto.OrderHistoryResponse{
			Id:            fmt.Sprintf("%s%04d", "7", item.Id),
			OrderDate:     item.OrderedAt.Format(time.RFC3339),
			TotalPayment:  item.TotalPayment.String(),
			Payment:       item.Payment.String(),
			PaymentReturn: item.PaymentReturn.String(),
			PaymentMethod: item.PaymentMethod,
		}
		for _, order := range data[i].OrderItems {
			orderItem := &dto.OrderItemResponse{
				Id:       order.ProductId,
				Name:     order.Product.Name,
				Quantity: order.Quantity,
				SubTotal: order.SubTotal.String(),
			}
			orders.OrderItem = append(orders.OrderItem, orderItem)
		}
		listOrders = append(listOrders, orders)
	}
	c.JSON(200, dto.Response{
		Data:        listOrders,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}
