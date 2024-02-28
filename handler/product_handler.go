package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/minimarket-app/dto"
	"github.com/ziadrahmatullah/minimarket-app/usecase"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(usecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: usecase}
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var request dto.AddProductReq
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	product, err := request.ToProduct()
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.usecase.AddProduct(c.Request.Context(), product)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, dto.Response{
		Message: "create success",
	})
}
