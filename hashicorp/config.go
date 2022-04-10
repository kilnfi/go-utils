package hashicorp

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"golang.org/x/time/rate"

	"github.com/skillz-blockchain/go-utils/common"
	kilntypes "github.com/skillz-blockchain/go-utils/common/types"
	kilntls "github.com/skillz-blockchain/go-utils/crypto/tls"
	kilnhttp "github.com/skillz-blockchain/go-utils/net/http"
)

// ClientConfig object that be converted into an api.Config later
type ClientConfig struct {
	Address string

	Mount string

	AgentAddress string

	Auth *AuthConfig

	HTTP *kilnhttp.ClientConfig

	MinRetryWait *kilntypes.Duration
	MaxRetryWait *kilntypes.Duration
	MaxRetries   *int

	RateLimit *RateLimitConfig
}

type AuthConfig struct {
	Token       string
	GitHubToken string
}

type RateLimitConfig struct {
	Rate  float64
	Burst int
}

func (cfg *ClientConfig) SetDefault() *ClientConfig {
	if cfg.MinRetryWait == nil {
		cfg.MinRetryWait = &kilntypes.Duration{Duration: 1000 * time.Millisecond}
	}

	if cfg.MaxRetryWait == nil {
		cfg.MaxRetryWait = &kilntypes.Duration{Duration: 1500 * time.Millisecond}
	}

	if cfg.MaxRetries == nil {
		cfg.MaxRetries = common.IntPtr(2)
	}

	if cfg.RateLimit == nil {
		cfg.RateLimit = &RateLimitConfig{}
	}

	if cfg.HTTP == nil {
		cfg.HTTP = &kilnhttp.ClientConfig{}
	}
	cfg.HTTP.SetDefault()

	if cfg.HTTP.Transport.TLS == nil {
		cfg.HTTP.Transport.TLS = &kilntls.Config{}
	}
	cfg.HTTP.Transport.TLS.MinVersion = "VersionTLS12"
	cfg.HTTP.Transport.EnableHTTP2 = true

	return cfg
}

func (cfg *ClientConfig) ToHashicorpConfig() (*api.Config, error) {
	client, err := kilnhttp.NewClient(cfg.HTTP)
	if err != nil {
		return nil, err
	}

	// Create Hashicorp Configuration
	config := &api.Config{
		Address:      cfg.Address,
		AgentAddress: cfg.AgentAddress,
		HttpClient:   client,
		Limiter:      rate.NewLimiter(rate.Limit(cfg.RateLimit.Rate), cfg.RateLimit.Burst),
	}

	config.MinRetryWait = cfg.MinRetryWait.Duration
	config.MinRetryWait = cfg.MinRetryWait.Duration
	config.MaxRetries = *cfg.MaxRetries

	// Ensure redirects are not automatically followed
	// Note that this is sane for the API client as it has its own
	// redirect handling logic (and thus also for command/meta),
	// but in e.g. http_test actual redirect handling is necessary
	config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// Returning this value causes the Go net library to not close the
		// response body and to nil out the error. Otherwise, retry clients may
		// try three times on every redirect because it sees an error from this
		// function (to prevent redirects) passing through to it.
		return http.ErrUseLastResponse
	}

	config.Backoff = retryablehttp.LinearJitterBackoff

	return config, nil
}
