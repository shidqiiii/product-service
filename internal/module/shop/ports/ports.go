package ports

import (
	"context"
	"product-service/internal/module/shop/entity"
)

type ShopService interface {
	CreateShop(ctx context.Context, req *entity.CreateShopRequest) (entity.UpsertShopResponse, error)
	DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error
	GetShops(ctx context.Context, req *entity.GetShopsRequest) (entity.GetShopsResponse, error)
	UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (entity.UpsertShopResponse, error)
}

type ShopRepository interface {
	CreateShop(ctx context.Context, req *entity.CreateShopRequest) (entity.UpsertShopResponse, error)
	DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error
	GetShops(ctx context.Context, req *entity.GetShopsRequest) (entity.GetShopsResponse, error)
	UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (entity.UpsertShopResponse, error)
	IsHaveShop(ctx context.Context, UserId string) (bool, error)
}
