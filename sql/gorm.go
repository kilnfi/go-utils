package sql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GormOpen(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(
		postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        cfg.DSN(),
		}),
		new(gorm.Config),
	)
}
