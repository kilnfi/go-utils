package app

import (
	"time"

	types "github.com/kilnfi/go-utils/common/types"
	kilnlog "github.com/kilnfi/go-utils/log"
	kilnnet "github.com/kilnfi/go-utils/net"
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

type Config struct {
	Logger       *kilnlog.Config
	Server       *kilnhttp.ServerConfig
	Healthz      *kilnhttp.ServerConfig
	StartTimeout *types.Duration
	StopTimeout  *types.Duration
}

func (cfg *Config) SetDefault() *Config {
	if cfg.Logger == nil {
		cfg.Logger = &kilnlog.Config{}
	}
	cfg.Logger.SetDefault()

	if cfg.Server == nil {
		cfg.Server = &kilnhttp.ServerConfig{}
	}
	cfg.Server.SetDefault()

	if cfg.Healthz == nil {
		cfg.Healthz = &kilnhttp.ServerConfig{}
	}
	if cfg.Healthz.Entrypoint == nil {
		cfg.Healthz.Entrypoint = &kilnnet.EntrypointConfig{}
	}
	if cfg.Healthz.Entrypoint.Address == "" {
		cfg.Healthz.Entrypoint.Address = ":8081"
	}
	cfg.Healthz.SetDefault()

	if cfg.StartTimeout == nil {
		cfg.StartTimeout = &types.Duration{Duration: 10 * time.Second}
	}

	if cfg.StopTimeout == nil {
		cfg.StopTimeout = &types.Duration{Duration: 10 * time.Second}
	}

	return cfg
}
