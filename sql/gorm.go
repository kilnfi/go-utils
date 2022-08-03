package sql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GormOpen(cfg *Config, gormLoggerOn bool) (*gorm.DB, error) {
	config := new(gorm.Config)
	if gormLoggerOn {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	return gorm.Open(
		postgres.Open(cfg.DSN().String()),
		config,
	)
}
