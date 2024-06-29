package rest

import (
	"product-service/internal/adapter"
	m "product-service/internal/middleware"
	"product-service/internal/module/shop/entity"
	"product-service/internal/module/shop/ports"
	"product-service/internal/module/shop/repository"
	"product-service/internal/module/shop/service"
	"product-service/pkg/errmsg"
	"product-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type shopHandler struct {
	service ports.ShopService
}

func NewShopHandler() *shopHandler {
	repo := repository.NewShopRepo(adapter.Adapters.ShopeefunProductPostgres)
	service := service.NewShopService(repo)

	return &shopHandler{
		service: service,
	}
}

func (h *shopHandler) Register(router fiber.Router) {

	router.Post("/shops", m.AuthQueryParams, h.createShop)
	router.Delete("/shops/:id", m.AuthQueryParams, h.deleteShop)
}

func (h *shopHandler) createShop(c *fiber.Ctx) error {
	var (
		req = &entity.CreateShopRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.UserId = c.Query("user_id")

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("service: Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	log.Debug().Interface("payload", req).Msg("service: CreateShop")

	resp, err := h.service.CreateShop(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))
}

func (h *shopHandler) deleteShop(c *fiber.Ctx) error {
	var (
		req = &entity.DeleteShopRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")
	req.UserId = c.Query("user_id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.DeleteShop(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
