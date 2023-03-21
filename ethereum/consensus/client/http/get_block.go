package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
)

// GetBlock returns block details for given block id.
func (c *Client) GetBlock(ctx context.Context, blockID string) (*bellatrix.SignedBeaconBlock, error) {
	rv, err := c.getBlock(ctx, blockID)
	if err != nil {
		c.logger.
			WithField("block", blockID).
			WithError(err).Errorf("GetBlock failed")
	}

	return rv, err
}

func (c *Client) getBlock(ctx context.Context, blockID string) (*bellatrix.SignedBeaconBlock, error) {
	req, err := newGetBlockRequest(ctx, blockID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlock", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlock", resp, "Failure sending request")
	}

	result, err := inspectGetBlockResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlock", resp, "Invalid response")
	}

	return result, nil
}

func newGetBlockRequest(ctx context.Context, blockID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"blockID": autorest.Encode("path", blockID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v2/beacon/blocks/{blockID}", pathParameters),
	).Prepare(newRequest(ctx))
}

type getBlockResponseMsg struct {
	Version string                       `json:"version"`
	Data    *bellatrix.SignedBeaconBlock `json:"data"`
}

func inspectGetBlockResponse(resp *http.Response) (*bellatrix.SignedBeaconBlock, error) {
	msg := new(getBlockResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
