package http

import (
	kilnnet "github.com/kilnfi/go-utils/net"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ServerFlagPrefixer struct {
	ep kilnnet.EntrypointFlagPrefixer
}

func NewServerFlagPrefixer(name, defaultAddr string) ServerFlagPrefixer {
	return ServerFlagPrefixer{
		ep: *kilnnet.NewEntrypointFlagPrefixer(name, defaultAddr),
	}
}

func (fl *ServerFlagPrefixer) Flags(v *viper.Viper, f *pflag.FlagSet) {
	fl.ep.Flags(v, f)
}

func (fl *ServerFlagPrefixer) ConfigFromViper(v *viper.Viper) *ServerConfig {
	return &ServerConfig{
		Entrypoint: fl.ep.ConfigFromViper(v),
	}
}
