package repository

import (
	"context"
	"database/sql"
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

func (s *shopRepo) DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error {
	query := `
		UPDATE
			shops
		SET
			deleted_at = NOW()
		WHERE
			user_id = $1
			AND id = $2
			AND deleted_at IS NULL
	`

	_, err := s.db.ExecContext(ctx, query, req.UserId, req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: Failed to delete shop")
		return err
	}

	return nil
}

func (s *shopRepo) GetShops(ctx context.Context, req *entity.GetShopsRequest) (entity.GetShopsResponse, error) {
	type (
		dao struct {
			TotalData int `db:"total_data"`
			entity.ShopItem
		}
	)
	var (
		arg  = make(map[string]any)
		res  = entity.GetShopsResponse{}
		data = make([]dao, 0)
	)
	res.Meta.Page = req.Page
	res.Meta.Limit = req.Limit
	res.Items = make([]entity.ShopItem, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			user_id,
			name,
			created_at,
			updated_at,
			deleted_at
		FROM
			shops
		WHERE
			1 = 1
	`

	if req.ShopName != "" {
		query += ` AND name ILIKE '%' || :name || '%'`
		arg["name"] = req.ShopName
	}

	query += `
		ORDER BY
			created_at DESC
		LIMIT :limit
		OFFSET :offset
	`
	arg["limit"] = req.Limit
	arg["offset"] = (req.Page - 1) * req.Limit

	nstmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: Failed to prepare query")
		return res, err
	}
	defer nstmt.Close()

	err = nstmt.SelectContext(ctx, &data, arg)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository: Failed to get shops")
		return res, err
	}

	for _, item := range data {
		res.Items = append(res.Items, item.ShopItem)

		res.Meta.TotalData = item.TotalData
	}

	res.Meta.CountTotalPage()
	return res, nil
}

func (s *shopRepo) UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (entity.UpsertShopResponse, error) {
	var (
		res = entity.UpsertShopResponse{}
	)

	query := `
		UPDATE
			shops
		SET
			name = $1,
			updated_at = NOW()
		WHERE
			user_id = $2
			AND id = $3
		RETURNING
			id, user_id, name, created_at, updated_at
	`

	err := s.db.QueryRowxContext(ctx, query, req.Name, req.UserId, req.Id).StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Any("payload", req).Msg("repository: Shop not found")
			return res, errmsg.NewCostumErrors(404, errmsg.WithMessage("Shop not found"))

		}
		log.Error().Err(err).Any("payload", req).Msg("repository: Failed to update shop")
		return res, err
	}

	return res, nil
}
