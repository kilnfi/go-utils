package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

// GetNodeVersion returns node's version contains informations about the node processing the request
func (c *Client) GetNodeVersion(ctx context.Context) (string, error) {
	rv, err := c.getNodeVersion(ctx)
	if err != nil {
		c.logger.WithError(err).Errorf("GetNodeVersion failed")
	}

	return rv, err
}

func (c *Client) getNodeVersion(ctx context.Context) (string, error) {
	req, err := newGetNodeVersionRequest(ctx)
	if err != nil {
		return "", autorest.NewErrorWithError(err, "eth2http.Client", "GetNodeVersion", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", autorest.NewErrorWithError(err, "eth2http.Client", "GetNodeVersion", resp, "Failure sending request")
	}

	result, err := inspectGetNodeVersionResponse(resp)
	if err != nil {
		return "", autorest.NewErrorWithError(err, "eth2http.Client", "GetNodeVersion", resp, "Invalid response")
	}

	return result, nil
}

func newGetNodeVersionRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("eth/v1/node/version"),
	).Prepare(newRequest(ctx))
}

type versionMsg struct {
	Version string `json:"version"`
}

type getNodeVersionResponseMsg struct {
	Data versionMsg `json:"data"`
}

func inspectGetNodeVersionResponse(resp *http.Response) (string, error) {
	msg := new(getNodeVersionResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return "", err
	}

	return msg.Data.Version, nil
}
