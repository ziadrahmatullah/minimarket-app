package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/handler"
	"github.com/ziadrahmatullah/minimarket-app/middleware"
)

type Handlers struct {
	User *handler.UserHandler
	Auth *handler.AuthHandler
}

func New(handlers Handlers) http.Handler {
	router := gin.New()

	router.NoRoute(routeNotFoundHandler)
	router.Use(gin.Recovery())
	router.Use(middleware.Timeout())
	router.Use(middleware.Logger())
	router.Use(middleware.Error())

	user := router.Group("/users")
	user.GET("", handlers.User.GetAllUser)

	auth := router.Group("/auth")
	auth.POST("", handlers.Auth.Register)

	return router
}

func routeNotFoundHandler(c *gin.Context) {
	var errRouteNotFound = errors.New("route not found")
	_ = c.Error(apperror.NewClientError(errRouteNotFound).NotFound())
}
