package net

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	types "github.com/kilnfi/go-utils/common/types"
	kilntls "github.com/kilnfi/go-utils/crypto/tls"
	"github.com/sirupsen/logrus"
)

type EntrypointConfig struct {
	Network string `json:"network"`
	Address string `json:"address"`

	KeepAlive *types.Duration

	TLSConfig *kilntls.Config
}

func (cfg *EntrypointConfig) SetDefault() *EntrypointConfig {
	if cfg.Address == "" {
		cfg.Address = "localhost:8080"
	}

	if cfg.Network == "" {
		cfg.Network = "tcp"
	}

	if cfg.KeepAlive == nil {
		cfg.KeepAlive = &types.Duration{Duration: 90 * time.Second}
	}

	return cfg
}

type Entrypoint struct {
	cfg *EntrypointConfig

	tlsCfg *tls.Config

	logger logrus.FieldLogger
}

func NewEntrypoint(cfg *EntrypointConfig) (*Entrypoint, error) {
	l := &Entrypoint{
		cfg: cfg,
	}

	if cfg.TLSConfig != nil {
		tlsCfg, err := cfg.TLSConfig.ToTLSConfig()
		if err != nil {
			return nil, err
		}
		l.tlsCfg = tlsCfg
	}

	return l, nil
}

func (lstnr *Entrypoint) SetLogger(logger logrus.FieldLogger) {
	lstnr.logger = logger.
		WithField("component", "entrypoint").
		WithField("network", lstnr.cfg.Network).
		WithField("address", lstnr.cfg.Address)
}

func (lstnr *Entrypoint) Logger() logrus.FieldLogger {
	if lstnr.logger == nil {
		lstnr.SetLogger(logrus.StandardLogger())
	}
	return lstnr.logger
}

func (lstnr *Entrypoint) Listen(ctx context.Context) (l net.Listener, err error) {
	lc := net.ListenConfig{
		KeepAlive: lstnr.cfg.KeepAlive.Duration,
	}

	logger := lstnr.Logger()

	logger.Infof("start listening for entering connection")
	l, err = lc.Listen(ctx, lstnr.cfg.Network, lstnr.cfg.Address)
	if err != nil {
		logger.WithError(err).Infof("error starting entrypoint")
		return
	}

	if lstnr.cfg.TLSConfig != nil {
		logger.Infof("upgrade to TLS")
		l = tls.NewListener(l, lstnr.tlsCfg)
	}

	return l, nil
}
