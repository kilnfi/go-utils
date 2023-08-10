//go:build integration
// +build integration

package docker

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func testService(t *testing.T, svcName string, svcCfg *ServiceConfig) {
	// Create compose
	cfg := (&ComposeConfig{
		Namespace: "test",
	}).SetDefault()
	compose, err := NewCompose(cfg)
	require.NoError(t, err)

	// Register service
	compose.RegisterService(svcName, svcCfg)

	// compose up
	err = compose.Up(context.TODO())
	require.NoError(t, err, "Up must not error")

	// wait for container to be ready
	err = compose.WaitContainer(context.TODO(), svcName, 5*time.Second)
	require.NoError(t, err, "WaitContainer must not error")

	// test container
	container, err := compose.GetContainer(context.TODO(), svcName)
	require.NoError(t, err, "GetContainer must not error")
	assert.Equal(t, fmt.Sprintf("/test_%v", svcName), container.Name)

	// compose down
	err = compose.Down(context.TODO())
	require.NoError(t, err, "Down must not error")
}

func TestPostgresService(t *testing.T) {
	svcCfg, err := NewPostgresServiceConfig(new(PostgresServiceOpts).SetDefault())
	require.NoError(t, err)

	testService(t, "postgres", svcCfg)
}

func TestTraefikService(t *testing.T) {
	svcCfg, err := NewTreafikServiceConfig(new(TraefikServiceOpts).SetDefault())
	require.NoError(t, err)

	testService(t, "traefik", svcCfg)
}

func TestFoundryService(t *testing.T) {
	svcCfg, err := NewFoundryServiceConfig(new(FoundryServiceOpts).SetDefault())
	require.NoError(t, err)

	testService(t, "foundry", svcCfg)
}
