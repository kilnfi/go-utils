package sql

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GormOpen(cfg *Config) (*gorm.DB, error) {
	config := new(gorm.Config)
	if cfg.GormLoggerOff {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	return gorm.Open(
		postgres.Open(cfg.DSN().String()),
		config,
	)
}

func PingGormDB(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return PingDB(ctx, sqlDB)
}
