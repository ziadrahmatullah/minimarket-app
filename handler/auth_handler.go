package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/minimarket-app/dto"
	"github.com/ziadrahmatullah/minimarket-app/usecase"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: u,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var request dto.RegisterReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		_ = ctx.Error(err)
		return
	}
	user := request.ToUser()
	err := h.usecase.Register(ctx.Request.Context(), user)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{Message: "register success"})
}


