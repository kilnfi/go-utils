//go:build !integration
// +build !integration

package sql

import (
	"testing"
	"time"

	types "github.com/kilnfi/go-utils/common/types"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	fl := NewFlagPrefixer("postgres", "Test")
	v := viper.New()

	fl.Flags(v, pflag.NewFlagSet("test", pflag.ContinueOnError))
	t.Setenv("TEST_DB_USER", "testuser")
	t.Setenv("TEST_DB_PASSWORD", "testpwd")
	t.Setenv("TEST_DB_HOST", "testhost")
	t.Setenv("TEST_DB_PORT", "1234")
	t.Setenv("TEST_DB_NAME", "testdb")
	t.Setenv("TEST_DB_SSLMODE", "verify-full")
	t.Setenv("TEST_DB_SSLCA", "./ca.pem")
	t.Setenv("TEST_DB_SSLCERT", "./cert.pem")
	t.Setenv("TEST_DB_SSLKEY", "./key.pem")
	t.Setenv("TEST_DB_CONNECT_TIMEOUT", "120s")

	cfg := fl.ConfigFromViper(v)
	assert.Equal(
		t,
		&Config{
			Dialect:        "postgres",
			User:           "testuser",
			Password:       "testpwd",
			Host:           "testhost",
			Port:           1234,
			DBName:         "testdb",
			SSLMode:        "verify-full",
			SSLCA:          "./ca.pem",
			SSLCert:        "./cert.pem",
			SSLKey:         "./key.pem",
			ConnectTimeout: &types.Duration{Duration: 120 * time.Second},
		},
		cfg,
	)
}
