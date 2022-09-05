//go:build integration
// +build integration

package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	kilndocker "github.com/kilnfi/go-utils/docker"
	kilnsql "github.com/kilnfi/go-utils/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareComposeDatabase(t *testing.T) *kilnsql.Config {
	// Create compose
	cfg := (&kilndocker.ComposeConfig{
		Namespace: "test.sql",
	}).SetDefault()
	compose, err := kilndocker.NewCompose(cfg)
	require.NoError(t, err)

	opts := new(kilndocker.PostgresServiceOpts).SetDefault()
	svcCfg, err := kilndocker.NewPostgresServiceConfig(opts)
	require.NoError(t, err)

	svcName := "postgres"
	compose.RegisterService(svcName, svcCfg)

	err = compose.Up(context.TODO())
	require.NoError(t, err, "Up must not error")

	t.Cleanup(func() {
		err = compose.Down(context.TODO())
		require.NoError(t, err, "Down must not error")
	})

	// wait for container to be ready
	err = compose.WaitContainer(context.TODO(), svcName, 5*time.Second)
	require.NoError(t, err, "WaitContainer must not error")

	container, err := compose.GetContainer(context.TODO(), svcName)
	require.NoError(t, err, "GetContainer must not error")

	sqlCfg, err := opts.SQLConfig(container)
	require.NoError(t, err, "GetContainer must not error")

	return sqlCfg
}

func TestCreateTempDB(t *testing.T) {
	sqlCfg := prepareComposeDatabase(t)

	cfg, err := kilnsql.CreateTempDB(t, sqlCfg)
	require.NoError(t, err)

	conn, err := pgx.Connect(context.TODO(), cfg.DSN().String())
	require.NoError(t, err)

	err = kilnsql.PingPGXConn(context.TODO(), conn)
	assert.NoError(t, err)
}
