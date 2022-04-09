package net

import (
	"net"
	"time"

	"github.com/skillz-blockchain/go-utils/common"
)

type DialerConfig struct {
	Timeout   *common.Duration
	KeepAlive *common.Duration
}

func (cfg *DialerConfig) SetDefault() *DialerConfig {
	if cfg.Timeout == nil {
		cfg.Timeout = &common.Duration{Duration: 30 * time.Second}
	}

	if cfg.KeepAlive == nil {
		cfg.KeepAlive = &common.Duration{Duration: 30 * time.Second}
	}

	return cfg
}

func NewDialer(cfg *DialerConfig) *net.Dialer {
	return &net.Dialer{
		Timeout:   cfg.Timeout.Duration,
		KeepAlive: cfg.KeepAlive.Duration,
	}
}
