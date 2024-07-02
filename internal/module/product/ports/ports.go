package ports

import (
	"context"
	"product-service/internal/module/product/entity"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (entity.UpsertProductResponse, error)
	GetProducts(ctx context.Context, req *entity.GetProductsRequest) (entity.GetProductsResponse, error)
	GetProduct(ctx context.Context, req *entity.GetProductRequest) (entity.GetProductResponse, error)
	UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (entity.UpsertProductResponse, error)
	DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (entity.UpsertProductResponse, error)
	GetProducts(ctx context.Context, req *entity.GetProductsRequest) (entity.GetProductsResponse, error)
	GetProduct(ctx context.Context, req *entity.GetProductRequest) (entity.GetProductResponse, error)
	UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (entity.UpsertProductResponse, error)
	DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error

	IsShopOwner(ctx context.Context, userId, shopId string) (bool, error)
	IsProductOwner(ctx context.Context, userId, productId string) (bool, error)
}
