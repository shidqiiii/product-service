package repository

import (
	"context"
	"product-service/internal/module/shop/entity"
	"product-service/internal/module/shop/ports"
	"product-service/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type shopRepo struct {
	db *sqlx.DB
}

func NewShopRepo(db *sqlx.DB) ports.ShopRepository {
	return &shopRepo{
		db: db,
	}
}

func (s *shopRepo) CreateShop(ctx context.Context, req *entity.CreateShopRequest) (entity.UpsertShopResponse, error) {
	var (
		res = entity.UpsertShopResponse{}
	)

	if exist, err := s.IsHaveShop(ctx, req.UserId); err != nil {
		return res, err
	} else if exist {
		log.Warn().Any("payload", req).Msg("repository: User already have shop")
		return res, errmsg.NewCostumErrors(400, errmsg.WithMessage("User already have shop"))
	}

	query := `
		INSERT INTO
			shops (user_id, name)
		VALUES ($1, $2)
		RETURNING
			id,
			user_id,
			name,
			created_at,
			updated_at
	`

	err := s.db.QueryRowContext(ctx, query, req.UserId, req.Name).Scan(
		&res.Id,
		&res.UserId,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: Failed to create shop")
		return res, err
	}

	return res, nil
}

func (s *shopRepo) IsHaveShop(ctx context.Context, UserId string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM shops WHERE user_id = $1 AND deleted_at IS NULL)`
	var exist bool
	err := s.db.GetContext(ctx, &exist, query, UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", UserId).Msg("repository: Failed to check shop existance")
		return false, err
	}

	return exist, nil
}
