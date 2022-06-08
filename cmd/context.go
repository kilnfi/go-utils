package cmd

import (
	"context"

	"github.com/spf13/viper"
)

type ctxKey string

var viperCtxKey ctxKey = "viper"

func WithViper(ctx context.Context, v *viper.Viper) context.Context {
	return context.WithValue(ctx, viperCtxKey, v)
}

func ViperFromContext(ctx context.Context) *viper.Viper {
	v, ok := ctx.Value(viperCtxKey).(*viper.Viper)
	if ok {
		return v
	}
	return viper.GetViper()
}
