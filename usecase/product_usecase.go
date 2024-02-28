package usecase

import (
	"context"

	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/repository"
	"github.com/ziadrahmatullah/minimarket-app/valueobject"
)

type ProductUsecase interface {
	AddProduct(ctx context.Context, product *entity.Product) error
	ListAllProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
}

type productUsecase struct {
	productRepo         repository.ProductRepository
	productCategoryRepo repository.ProductCategoryRepository
}

func NewProductUsecase(
	productRepo repository.ProductRepository,
	productCategoryRepo repository.ProductCategoryRepository,
) ProductUsecase {
	return &productUsecase{
		productRepo:         productRepo,
		productCategoryRepo: productCategoryRepo,
	}
}

func (u *productUsecase) AddProduct(ctx context.Context, product *entity.Product) error {
	fetchedProductCategory, err := u.productCategoryRepo.FindById(ctx, product.ProductCategoryId)
	if err != nil {
		return err
	}
	if fetchedProductCategory == nil {
		return apperror.NewResourceNotFoundError("product category", "id", product.ProductCategoryId)
	}
	productQuery := valueobject.NewQuery().Condition("product_code", valueobject.Equal, product.ProductCode)
	fetchedProduct, err := u.productRepo.FindOne(ctx, productQuery)
	if err != nil {
		return err
	}
	if fetchedProduct != nil {
		return apperror.NewResourceAlreadyExistError("product", "produc_code", product.Name)
	}
	_, err = u.productRepo.Create(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (u *productUsecase) ListAllProduct(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	query.WithJoin("ProductCategory")
	pagedResult, err := u.productRepo.FindAllProducts(ctx, query)
	if err != nil {
		return nil, err
	}

	return pagedResult, nil
}
