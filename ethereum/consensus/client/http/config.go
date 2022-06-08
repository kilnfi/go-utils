package eth2http

import (
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

type Config struct {
	Address string

	HTTP *kilnhttp.ClientConfig
}

func (cfg *Config) SetDefault() *Config {
	if cfg.HTTP == nil {
		cfg.HTTP = new(kilnhttp.ClientConfig)
	}

	cfg.HTTP.SetDefault()

	return cfg
}
