package sql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GormOpen(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(
		postgres.Open(cfg.DSN().String()),
		new(gorm.Config),
	)
}
