//go:build !integration
// +build !integration

package utils

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	assert.Equal(t, viper.GetViper(), ViperFromContext(context.TODO()))

	newV := viper.New()
	ctx := WithViper(context.Background(), newV)

	newV.Set("test", "test")
	assert.Equal(t, newV, ViperFromContext(ctx))
}
