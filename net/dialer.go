package net

import (
	"net"
	"time"

	kilntypes "github.com/kilnfi/go-utils/common/types"
)

type DialerConfig struct {
	Timeout   *kilntypes.Duration
	KeepAlive *kilntypes.Duration
}

func (cfg *DialerConfig) SetDefault() *DialerConfig {
	if cfg.Timeout == nil {
		cfg.Timeout = &kilntypes.Duration{Duration: 30 * time.Second}
	}

	if cfg.KeepAlive == nil {
		cfg.KeepAlive = &kilntypes.Duration{Duration: 30 * time.Second}
	}

	return cfg
}

func NewDialer(cfg *DialerConfig) *net.Dialer {
	return &net.Dialer{
		Timeout:   cfg.Timeout.Duration,
		KeepAlive: cfg.KeepAlive.Duration,
	}
}
