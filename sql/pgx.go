package sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func PGXConnect(ctx context.Context, cfg *Config) (*pgx.Conn, error) {
	return pgx.Connect(context.TODO(), cfg.DSN().String())
}

func PingPGXConn(ctx context.Context, conn *pgx.Conn) error {
	err := conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("health check failed on ping: %w", err)
	}

	rows, err := conn.Query(ctx, `SELECT VERSION()`)
	if err != nil {
		return fmt.Errorf("health check failed on select: %w", err)
	}

	rows.Close()

	return nil
}
