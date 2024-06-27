package ports

import (
	"context"
	"product-service/internal/module/shop/entity"
)

type ShopService interface {
	CreateShop(ctx context.Context, req *entity.CreateShopRequest) (entity.UpsertShopResponse, error)
}

type ShopRepository interface {
	CreateShop(ctx context.Context, req *entity.CreateShopRequest) (entity.UpsertShopResponse, error)
}
