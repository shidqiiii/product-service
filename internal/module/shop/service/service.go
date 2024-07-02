package service

import (
	"context"
	"product-service/internal/module/shop/entity"
	"product-service/internal/module/shop/ports"
	"product-service/pkg/errmsg"

	"github.com/rs/zerolog/log"
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
	var (
		res = entity.UpsertShopResponse{}
	)

	exist, err := s.repo.IsHaveShop(ctx, req.UserId)
	if err != nil {
		return res, err

	}

	if exist {
		log.Warn().Any("payload", req).Msg("service: Creating shop")
		return res, errmsg.NewCostumErrors(400, errmsg.WithMessage("User already have shop"))
	}

	result, err := s.repo.CreateShop(ctx, req)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *shopService) DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error {
	return s.repo.DeleteShop(ctx, req)
}

func (s *shopService) GetShops(ctx context.Context, req *entity.GetShopsRequest) (entity.GetShopsResponse, error) {
	return s.repo.GetShops(ctx, req)
}

func (s *shopService) UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (entity.UpsertShopResponse, error) {
	return s.repo.UpdateShop(ctx, req)
}
