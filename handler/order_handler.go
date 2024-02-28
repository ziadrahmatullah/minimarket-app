package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/minimarket-app/dto"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/usecase"
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

func (h *OrderHandler) DailyOrderReport(c *gin.Context){
	var requestParam dto.ReportDailyQueryParamReq
	if err := c.ShouldBindQuery(&requestParam); err != nil {
		_ = c.Error(err)
		return
	}
	query := requestParam.ToQuery()
	pageResult, err := h.usecase.DailyOrderReport(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	orders := pageResult.Data.([]*entity.Order)
	c.JSON(http.StatusOK, dto.Response{
		Data:        orders,
		TotalPage:   &pageResult.TotalPage,
		TotalItem:   &pageResult.TotalItem,
		CurrentPage: &pageResult.CurrentPage,
		CurrentItem: &pageResult.CurrentItems,
	})
}