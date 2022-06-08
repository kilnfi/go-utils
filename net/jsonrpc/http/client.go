package jsonrpchttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"

	kilnhttp "github.com/kilnfi/go-utils/net/http"
	httppreparer "github.com/kilnfi/go-utils/net/http/preparer"
	"github.com/kilnfi/go-utils/net/jsonrpc"
)

// Client allows to connect to a JSON-RPC server
type Client struct {
	client autorest.Sender

	logger logrus.FieldLogger
}

// NewClient creates a new client connected to a JSON-RPC server
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
	c.logger = logger.WithField("component", "jsonrpc.http-client")
}

func (c *Client) Call(ctx context.Context, r *jsonrpc.Request, res interface{}) error {
	err := c.call(ctx, r, res)
	if err != nil {
		c.logger.
			WithField("req.method", r.Method).
			WithField("req.version", r.Version).
			WithField("req.params", r.Params).
			WithField("req.id", r.ID).
			WithError(err).Errorf("jsonrpc call failed")
	}

	return err
}

// Call performs JSON-RPC call
func (c *Client) call(ctx context.Context, r *jsonrpc.Request, res interface{}) error {
	req, err := newCallRequest(ctx, r)
	if err != nil {
		return autorest.NewErrorWithError(err, "jsonrpchttp.Client", "Call", nil, "Request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		msg, _ := json.Marshal(r)
		return autorest.NewErrorWithError(err, "jsonrpchttp.Client", fmt.Sprintf("Call(%v)", string(msg)), resp, "Do")
	}

	err = inspectCallResponse(resp, res)
	if err != nil {
		msg, _ := json.Marshal(r)
		return autorest.NewErrorWithError(err, "jsonrpchttp.Client", fmt.Sprintf("Call(%v)", string(msg)), resp, "Response")
	}

	return nil
}

// ByUnmarshallingResponse marshall JSON-RPC request message into http.Request body
func newCallRequest(ctx context.Context, req *jsonrpc.Request) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithPath("/"),
		autorest.AsJSON(),
		autorest.WithJSON(req),
	).Prepare(newRequest(ctx))
}

func newRequest(ctx context.Context) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, "", "", http.NoBody)
	return req
}

// responseMsg is a struct allowing to encode/decode a JSON-RPC response body
type responseMsg struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *json.RawMessage `json:"error,omitempty"`
	ID      *json.RawMessage `json:"id,omitempty"`
}

func inspectCallResponseMsg(msg *responseMsg, res interface{}) error {
	if msg.Error == nil && msg.Result == nil {
		return fmt.Errorf("invalid JSON-RPC response missing both result and error")
	}

	if msg.Error != nil {
		errMsg := new(jsonrpc.ErrorMsg)
		err := json.Unmarshal(*msg.Error, errMsg)
		if err != nil {
			return fmt.Errorf("invalid JSON-RPC error message %v", string(*msg.Error))
		}
		return errMsg
	}

	if msg.Result != nil && res != nil {
		err := json.Unmarshal(*msg.Result, res)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON-RPC result %v into %T (%v)", string(*msg.Result), res, err)
		}
		return nil
	}

	return nil

}

func inspectCallResponse(resp *http.Response, res interface{}) error {
	msg := new(responseMsg)
	err := autorest.Respond(
		resp,
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(msg),
		autorest.ByClosing(),
	)
	if err != nil {
		return err
	}

	return inspectCallResponseMsg(msg, res)
}
