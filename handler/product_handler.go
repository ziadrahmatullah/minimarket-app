package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/minimarket-app/dto"
	"github.com/ziadrahmatullah/minimarket-app/entity"
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

func (h *ProductHandler) ListProduct(c *gin.Context) {
	var request dto.ListProductQueryParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}

	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}

	pagedResult, err := h.usecase.ListAllProduct(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	products := pagedResult.Data.([]*entity.Product)

	var response []*dto.ProductRes
	for _, product := range products {
		response = append(response, dto.NewFromProduct(product))
	}
	c.JSON(200, dto.Response{
		Data:        response,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}
