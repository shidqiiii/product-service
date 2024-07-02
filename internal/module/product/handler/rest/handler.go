package rest

import (
	"product-service/internal/adapter"
	m "product-service/internal/middleware"
	"product-service/internal/module/product/entity"
	"product-service/internal/module/product/ports"
	"product-service/internal/module/product/repository"
	"product-service/internal/module/product/service"
	"product-service/pkg/errmsg"
	"product-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type producthandler struct {
	service ports.ProductService
}

func NewProductHandler() *producthandler {
	repo := repository.NewProductRepository(adapter.Adapters.ShopeefunProductPostgres)
	service := service.NewProductService(repo)

	return &producthandler{
		service: service,
	}
}

func (h *producthandler) Register(router fiber.Router) {
	router.Get("/products", h.getProducts)
	router.Get("/products/:id", h.getProduct)

	router.Post("/products", m.AuthQueryParams, h.createProduct)
	router.Patch("/products/:id", m.AuthQueryParams, h.updateProduct)
	router.Delete("/products/:id", m.AuthQueryParams, h.deleteProduct)
}

func (h *producthandler) createProduct(c *fiber.Ctx) error {
	var (
		req = &entity.CreateProductRequest{}
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

	resp, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}

func (h *producthandler) updateProduct(c *fiber.Ctx) error {
	var (
		req = &entity.UpdateProductRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.UserId = c.Query("user_id")
	req.Id = c.Params("id")

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("service: Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.UpdateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *producthandler) deleteProduct(c *fiber.Ctx) error {
	var (
		req = &entity.DeleteProductRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.UserId = c.Query("user_id")
	req.ProductId = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.DeleteProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *producthandler) getProducts(c *fiber.Ctx) error {
	var (
		req = &entity.GetProductsRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Error().Err(err).Msg("service: Failed to parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefaults()

	if code, errs := req.CostumValidation(); code != 0 {
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request query")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProducts(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *producthandler) getProduct(c *fiber.Ctx) error {
	var (
		req = &entity.GetProductRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.ProductId = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("service: Invalid request query")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
