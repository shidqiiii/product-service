package infrastructure

import (
	"product-service/pkg/config"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	Envs *Config // Envs is global vars Config.
	once sync.Once
)

type Config struct {
	App struct {
		Environtment string `env:"APP_ENV" env-default:"development"`
		Port         string `env:"APP_PORT" env-default:"3000"`
		LogLevel     string `env:"APP_LOG_LEVEL" env-default:"debug"`
	}
	DB struct {
		ConnectionTimeout int `env:"DB_CONN_TIMEOUT" env-default:"30" env-description:"database timeout in seconds"`
		MaxOpenCons       int `env:"DB_MAX_OPEN_CONS" env-default:"20" env-description:"database max open conn in seconds"`
		MaxIdleCons       int `env:"DB_MAX_IdLE_CONS" env-default:"20" env-description:"database max idle conn in seconds"`
		ConnMaxLifetime   int `env:"DB_CONN_MAX_LIFETIME" env-default:"0" env-description:"database conn max lifetime in seconds"`
	}
	Guard struct {
		JwtPrivateKey string `env:"JWT_PRIVATE_KEY"`
	}
	ShopeefunProductPostgres struct {
		Host     string `env:"SHOPEEFUN_PRODUCT_POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"SHOPEEFUN_PRODUCT_POSTGRES_PORT" env-default:"5432"`
		Username string `env:"SHOPEEFUN_PRODUCT_POSTGRES_USER" env-default:"postgres"`
		Password string `env:"SHOPEEFUN_PRODUCT_POSTGRES_PASSWORD" env-default:"postgres"`
		Database string `env:"SHOPEEFUN_PRODUCT_POSTGRES_DB" env-default:"venatronics"`
		SslMode  string `env:"SHOPEEFUN_PRODUCT_POSTGRES_SSL_MODE" env-default:"disable"`
	}
}

// Option is Configure type return func.
type Option = func(c *Configure) error

// Configure is the data struct.
type Configure struct {
	path     string
	filename string
}

// Configuration create instance.
func Configuration(opts ...Option) *Configure {
	c := &Configure{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err)
		}
	}
	return c
}

// Initialize will create instance of Configure.
func (c *Configure) Initialize() {
	once.Do(func() {
		Envs = &Config{}
		if err := config.Load(config.Opts{
			Config:    Envs,
			Paths:     []string{c.path},
			Filenames: []string{c.filename},
		}); err != nil {
			log.Fatal().Err(err).Msg("get config error")
		}
	})
}

// WithPath will assign to field path Configure.
func WithPath(path string) Option {
	return func(c *Configure) error {
		c.path = path
		return nil
	}
}

// WithFilename will assign to field name Configure.
func WithFilename(name string) Option {
	return func(c *Configure) error {
		c.filename = name
		return nil
	}
}
