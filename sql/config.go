package sql

import (
	"fmt"
	"net/url"
	"time"

	types "github.com/kilnfi/go-utils/common/types"
)

type Config struct {
	Dialect  string
	User     string
	Password string
	Host     string
	Port     uint16
	DBName   string

	SSLMode, SSLCert, SSLKey, SSLCA string

	ConnectTimeout *types.Duration

	PoolSize int

	KeepAlive *types.Duration
}

var pgDialect = "postgres"

func (cfg *Config) SetDefault() *Config {
	if cfg.Dialect == "" {
		cfg.Dialect = pgDialect
	}

	if cfg.User == "" {
		cfg.User = pgDialect
	}

	if cfg.Password == "" {
		cfg.Password = pgDialect
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}

	if cfg.Port == 0 {
		cfg.Port = uint16(5432)
	}

	if cfg.DBName == "" {
		cfg.DBName = pgDialect
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	if cfg.ConnectTimeout == nil {
		cfg.ConnectTimeout = &types.Duration{Duration: 30 * time.Second}
	}

	if cfg.KeepAlive == nil {
		cfg.KeepAlive = &types.Duration{Duration: time.Minute}
	}

	return cfg
}

func (cfg *Config) DSN() string {
	u := url.URL{
		Scheme: cfg.Dialect,
		Host:   fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
		User:   url.UserPassword(cfg.User, cfg.Password),
		Path:   cfg.DBName,
	}

	query := make(url.Values)
	if cfg.ConnectTimeout != nil {
		query.Add("connect_timeout", fmt.Sprintf("%v", int64(cfg.ConnectTimeout.Duration/time.Second)))
	}

	if cfg.SSLCA != "" {
		query.Add("sslca", cfg.SSLCA)
	}

	if cfg.SSLKey != "" {
		query.Add("sslkey", cfg.SSLKey)
	}

	if cfg.SSLCert != "" {
		query.Add("sslcert", cfg.SSLCert)
	}

	if cfg.SSLMode != "" {
		query.Add("sslmode", cfg.SSLMode)
	}

	u.RawQuery = query.Encode()

	return u.String()
}
