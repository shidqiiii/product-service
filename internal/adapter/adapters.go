package adapter

import (
	"fmt"
	"strings"

	// import "product-service/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var (
	Adapters *Adapter
)

type Option func(adapter *Adapter)

type Validator interface {
	Validate(i any) error
}

type Adapter struct {
	// Driver Adapters
	RestServer *fiber.App

	// Driven Adapters
	ShopeefunProductPostgres *sqlx.DB
	Validator                Validator // *validator.Validator
}

func (a *Adapter) Sync(opts ...Option) {
	for o := range opts {
		opt := opts[o]
		opt(a)
	}
}

func (a *Adapter) Unsync() error {
	var errs []string

	if a.RestServer != nil {
		if err := a.RestServer.Shutdown(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Rest server disconnected")
	}

	if a.ShopeefunProductPostgres != nil {
		if err := a.ShopeefunProductPostgres.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Crowners Postgres disconnected")
	}

	if len(errs) > 0 {
		err := fmt.Errorf(strings.Join(errs, "\n"))
		log.Error().Msgf("Error while disconnecting adapters: %v", err)
		return err
	}

	return nil
}
