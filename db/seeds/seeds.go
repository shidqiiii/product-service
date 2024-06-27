package seeds

import (
	"os"
	"product-service/internal/adapter"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// Run seeds.
func (s *Seed) run(table string, total int) {

	switch table {
	case "users":
		s.usersSeed(total)
	case "product_categories":
		s.ProductCategoriesSeed(total)
	default:
		log.Info().Msg("No seed to run")
	}

	if table != "" {
		log.Info().Msg("Seed ran successfully")
		log.Info().Msg("Exiting ...")
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
		os.Exit(0)
	}

}

// this function is used to seed the product_categories table
// with random data
func (s *Seed) ProductCategoriesSeed(total int) {
	var (
		args  = make([]map[string]any, 0)
		query = "INSERT INTO product_categories (name) VALUES (:name)"
	)

	for i := 0; i < total; i++ {
		var (
			name = gofakeit.ProductCategory()
			arg  = make(map[string]any)
		)

		arg["name"] = name
		args = append(args, arg)
	}

	_, err := s.db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Error creating product categories")
	}

	log.Info().Msg("product_categories table seeded successfully")
}

// this function is used to seed the users table
// with random data
func (s *Seed) usersSeed(total int) {
	var (
		roles = []string{"admin", "buyer", "seller"}
		args  = make([]map[string]any, 0)
		query = `
			INSERT INTO users (username, email, role, address)
			VALUES (:username, :email, :role, :address)
		`
	)

	for i := 0; i < total; i++ {
		var (
			username  = gofakeit.Username()
			email     = gofakeit.Email()
			role      = roles[gofakeit.Number(0, 2)]
			addresses = []*string{&gofakeit.Address().Address, nil}
			address   = addresses[gofakeit.Number(0, 1)]
			arg       = make(map[string]any)
		)

		/*
			log.Info().Msg("Seeding user: " + username)
			log.Info().Msg("Email: " + email)
			log.Info().Msg("Role: " + role)
			if address != nil {
				log.Info().Msg("Address: " + *address)
			} else {
				log.Info().Msg("Address: nil")
			log.Info().Msg("====================================")
		*/

		arg["username"] = username
		arg["email"] = email
		arg["role"] = role
		arg["address"] = address
		args = append(args, arg)
	}

	_, err := s.db.NamedExec(query, args)
	if err != nil {
		log.Error().Err(err).Msg("Error creating users")
	}

	log.Info().Msg("users table seeded successfully")
}
