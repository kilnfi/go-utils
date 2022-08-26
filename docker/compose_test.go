//go:build integration
// +build integration

package docker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestCompose(t *testing.T) {
	cfg := (&ComposeConfig{
		Namespace: "test",
	}).SetDefault()
	compose, err := NewCompose(cfg)
	require.NoError(t, err)

	pgSvcCfg, err := NewPostgresServiceConfig(new(PostgresServiceOpts).SetDefault())
	require.NoError(t, err)

	svcName := "postgres"
	compose.RegisterService(svcName, pgSvcCfg)

	err = compose.Up(context.TODO())
	require.NoError(t, err, "Up must not error")

	err = compose.WaitContainer(context.TODO(), svcName, 5*time.Second)
	require.NoError(t, err, "WaitContainer must not error")

	container, err := compose.GetContainer(context.TODO(), svcName)
	require.NoError(t, err, "GetContainer must not error")

	assert.Equal(t, "/test_postgres", container.Name)

	err = compose.Down(context.TODO())
	require.NoError(t, err, "Down must not error")
}
