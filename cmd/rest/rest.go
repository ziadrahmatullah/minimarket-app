package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/ziadrahmatullah/minimarket-app/appjwt"
	"github.com/ziadrahmatullah/minimarket-app/appvalidator"
	"github.com/ziadrahmatullah/minimarket-app/handler"
	"github.com/ziadrahmatullah/minimarket-app/hasher"
	"github.com/ziadrahmatullah/minimarket-app/logger"
	"github.com/ziadrahmatullah/minimarket-app/repository"
	"github.com/ziadrahmatullah/minimarket-app/router"
	"github.com/ziadrahmatullah/minimarket-app/server"
	"github.com/ziadrahmatullah/minimarket-app/transactor"
	"github.com/ziadrahmatullah/minimarket-app/usecase"
)

func main() {
	logger.SetLogrusLogger()

	db, err := repository.GetConnection()
	if err != nil {
		logger.Log.Error(err)
	}

	client := resty.New()
	client.SetHeader("Content-Type", "application/json")

	manager := transactor.NewManager(db)
	hash := hasher.NewHasher()
	jwt := appjwt.NewJwt()
	appvalidator.RegisterCustomValidator()
	userR := repository.NewUserRepository(db)
	userU := usecase.NewUserUsecase(userR)
	userH := handler.NewUserHandler(userU)

	authU := usecase.NewAuthUsecase(userR, jwt, hash)
	authH := handler.NewAuthHandler(authU)

	productCategoryR := repository.NewProductCategoryRepository(db)
	productR := repository.NewProductRepository(db)
	productU := usecase.NewProductUsecase(productR, productCategoryR)
	productH := handler.NewProductHandler(productU)

	orderR := repository.NewOrderRepository(db)
	orderItemR := repository.NewOrderItemRepository(db)
	orderU := usecase.NewOrderUsecase(orderR, orderItemR, productR, manager)
	orderH := handler.NewOrderHandler(orderU)

	handlers := router.Handlers{
		User:    userH,
		Auth:    authH,
		Product: productH,
		Order:   orderH,
	}

	r := router.New(handlers)

	s := server.New(r)

	server.StartWithGracefulShutdown(s)
}
