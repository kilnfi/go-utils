package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"

	kilnhttp "github.com/skillz-blockchain/go-utils/net/http"
	httppreparer "github.com/skillz-blockchain/go-utils/net/http/preparer"
)

// Client provides methods to connect to an Ethereum 2.0 Beacon chain node
type Client struct {
	client autorest.Sender

	logger logrus.FieldLogger
}

func NewClientFromClient(s autorest.Sender) *Client {
	c := &Client{
		client: s,
	}

	c.SetLogger(logrus.StandardLogger())

	return c
}

// NewClient creates a client connecting to an Ethereum 2.0 Beacon chain node at given addr
func NewClient(cfg *Config) (*Client, error) {
	httpc, err := kilnhttp.NewClient(cfg.HTTP)
	if err != nil {
		return nil, err
	}

	return NewClientFromClient(
		autorest.Client{
			Sender:           httpc,
			RequestInspector: httppreparer.WithBaseURL(cfg.Address),
		},
	), nil
}

func (c *Client) Logger() logrus.FieldLogger {
	return c.logger
}

func (c *Client) SetLogger(logger logrus.FieldLogger) {
	c.logger = logger.WithField("component", "eth.consensus.client")
}

func newRequest(ctx context.Context) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, "", "", http.NoBody)
	return req
}

func inspectResponse(resp *http.Response, msg interface{}) error {
	return autorest.Respond(
		resp,
		WithBeaconErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(msg),
		autorest.ByClosing(),
	)
}
