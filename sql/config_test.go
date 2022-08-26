//go:build !integration
// +build !integration

package sql

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	types "github.com/kilnfi/go-utils/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	rootPEM = `-----BEGIN CERTIFICATE-----
MIIB0zCCAX2gAwIBAgIJAI/M7BYjwB+uMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTIwOTEyMjE1MjAyWhcNMTUwOTEyMjE1MjAyWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANLJ
hPHhITqQbPklG3ibCVxwGMRfp/v4XqhfdQHdcVfHap6NQ5Wok/4xIA+ui35/MmNa
rtNuC+BdZ1tMuVCPFZcCAwEAAaNQME4wHQYDVR0OBBYEFJvKs8RfJaXTH08W+SGv
zQyKn0H8MB8GA1UdIwQYMBaAFJvKs8RfJaXTH08W+SGvzQyKn0H8MAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADQQBJlffJHybjDGxRMqaRmDhX0+6v02TUKZsW
r5QuVbpQhH6u+0UgcW0jp9QwpxoPTLTWGXEWBBBurxFwiCBhkQ+V
-----END CERTIFICATE-----
`
)

func TestDSN(t *testing.T) {
	cfg := new(Config).SetDefault()

	assert.Equal(
		t,
		"postgres://postgres:postgres@localhost:5432/postgres?connect_timeout=30&sslmode=disable",
		cfg.DSN().String(),
	)

	cfg = &Config{
		Dialect:        "postgres",
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
		"postgres://user:pwd@host:1234/test?connect_timeout=120&sslcert=cert.pem&sslkey=key.pem&sslmode=require&sslrootcert=ca.pem",
		cfg.DSN().String(),
	)

	dir := t.TempDir()

	err := os.WriteFile(path.Join(dir, "ca.pem"), []byte(rootPEM), 0o777)
	require.NoError(t, err)

	cfg = &Config{
		Dialect:  "postgres",
		User:     "user",
		Password: "pwd",
		Host:     "host",
		Port:     1234,
		DBName:   "test",
		SSLMode:  "require",
		// SSLCert:        "cert.pem",
		// SSLKey:         "key.pem",
		SSLCA:          path.Join(dir, "ca.pem"),
		ConnectTimeout: &types.Duration{Duration: 120 * time.Second},
	}

	pgxCfg, err := pgconn.ParseConfig(cfg.DSN().String())
	require.NoError(t, err)

	assert.Equal(
		t,
		"user",
		pgxCfg.User,
	)
	assert.Equal(
		t,
		"pwd",
		pgxCfg.Password,
	)
	assert.Equal(
		t,
		"host",
		pgxCfg.Host,
	)
	assert.Equal(
		t,
		uint16(1234),
		pgxCfg.Port,
	)
	assert.Equal(
		t,
		"test",
		pgxCfg.Database,
	)
	assert.NotNil(
		t,
		pgxCfg.TLSConfig,
	)
	assert.Equal(t, true, pgxCfg.TLSConfig.InsecureSkipVerify)
	assert.Len(t, pgxCfg.TLSConfig.Certificates, 0)
	require.NotNil(t, pgxCfg.TLSConfig.RootCAs)
}
