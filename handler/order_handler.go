package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/minimarket-app/dto"
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
		TotalPayment: order.TotalPayment.String(), 
		PaymentReturn: order.PaymentReturn.String()}})
}

func (h *OrderHandler) GetMostOrderedCategories(c *gin.Context){
	categories, err := h.usecase.GetMostOrderedCategories(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: categories})
}