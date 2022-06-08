package http

import (
	"net/http"

	kilntypes "github.com/kilnfi/go-utils/common/types"
)

// Config for creating an HTTP Client
type ClientConfig struct {
	Transport *TransportConfig    `json:"transport,omitempty"`
	Timeout   *kilntypes.Duration `json:"timeout,omitempty"`
}

func (cfg *ClientConfig) SetDefault() *ClientConfig {
	if cfg.Transport == nil {
		cfg.Transport = new(TransportConfig)
	}
	cfg.Transport.SetDefault()

	if cfg.Timeout == nil {
		cfg.Timeout = &kilntypes.Duration{Duration: 0}
	}

	return cfg
}

// New creates a new HTTP client
func NewClient(cfg *ClientConfig) (*http.Client, error) {
	trnsprt, err := NewTransport(cfg.Transport)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: trnsprt,
		Timeout:   cfg.Timeout.Duration,
	}, nil
}
