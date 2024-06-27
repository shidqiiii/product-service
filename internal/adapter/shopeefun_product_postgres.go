package adapter

import (
	// "log"

	"product-service/internal/infrastructure"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func WithShopeefunProductPostgres() Option {
	return func(a *Adapter) {
		dbUser := infrastructure.Envs.ShopeefunProductPostgres.Username
		dbPassword := infrastructure.Envs.ShopeefunProductPostgres.Password
		dbName := infrastructure.Envs.ShopeefunProductPostgres.Database
		dbHost := infrastructure.Envs.ShopeefunProductPostgres.Host
		dbSSLMode := infrastructure.Envs.ShopeefunProductPostgres.SslMode
		dbPort := infrastructure.Envs.ShopeefunProductPostgres.Port

		dbMaxPoolSize := infrastructure.Envs.DB.MaxOpenCons
		dbMaxIdleConns := infrastructure.Envs.DB.MaxIdleCons
		dbConnMaxLifetime := infrastructure.Envs.DB.ConnMaxLifetime

		connectionString := "user=" + dbUser + " password=" + dbPassword + " host=" + dbHost + " port=" + dbPort + " dbname=" + dbName + " sslmode=" + dbSSLMode + " TimeZone=UTC"
		db, err := sqlx.Connect("postgres", connectionString)
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Postgres")
		}

		db.SetMaxOpenConns(dbMaxPoolSize)
		db.SetMaxIdleConns(dbMaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConnMaxLifetime) * time.Second)

		// check connection
		err = db.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Shopeefun Product Postgres")
		}

		a.ShopeefunProductPostgres = db
		log.Info().Msg("Shopeefun Product Postgres connected")
	}
}
