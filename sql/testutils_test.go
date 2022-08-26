package sql

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTempDB(t *testing.T) {
	cfg, err := CreateTempDB(t, (&Config{}).SetDefault())
	require.NoError(t, err)

	conn, err := pgx.Connect(context.TODO(), cfg.DSN().String())
	require.NoError(t, err)

	err = PingPGXConn(context.TODO(), conn)
	assert.NoError(t, err)
}
