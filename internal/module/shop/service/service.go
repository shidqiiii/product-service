package service

import (
	"context"
	"product-service/internal/module/shop/entity"
	"product-service/internal/module/shop/ports"
)

type shopService struct {
	repo ports.ShopRepository
}

func NewShopService(r ports.ShopRepository) ports.ShopService {
	return &shopService{
		repo: r,
	}
}

func (s *shopService) CreateShop(ctx context.Context, req *entity.CreateShopRequest) (entity.UpsertShopResponse, error) {
	result, err := s.repo.CreateShop(ctx, req)
	if err != nil {
		return result, err
	}

	return result, nil
}
