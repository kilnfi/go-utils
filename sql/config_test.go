package sql

import (
	"testing"
	"time"

	types "github.com/kilnfi/go-utils/common/types"
	"github.com/stretchr/testify/assert"
)

func TestDSN(t *testing.T) {
	cfg := new(Config).SetDefault()

	assert.Equal(
		t,
		"postgres://postgres:postgres@localhost:5432/postgres?connect_timeout=30&sslmode=disable",
		cfg.DSN(),
	)

	cfg = &Config{
		Dialect:        "mysql",
		User:           "user",
		Password:       "pwd",
		Host:           "host",
		Port:           1234,
		DBName:         "test",
		SSLMode:        "require",
		SSLCert:        "cert.pem",
		SSLKey:         "key.pem",
		SSLCA:          "ca.pem",
		ConnectTimeout: &types.Duration{Duration: 120 * time.Second},
	}

	assert.Equal(
		t,
		"mysql://user:pwd@host:1234/test?connect_timeout=120&sslca=ca.pem&sslcert=cert.pem&sslkey=key.pem&sslmode=require",
		cfg.DSN(),
	)
}
