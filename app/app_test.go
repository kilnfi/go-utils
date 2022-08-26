//go:build !integration
// +build !integration

package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	cfg := (&Config{}).SetDefault()

	_, err := New(cfg)

	require.NoError(t, err)
}
