package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"

	"github.com/kilnfi/go-utils/ethereum/consensus/types"
)

// GetBlockHeader returns block header for given blockID
func (c *Client) GetBlockHeader(ctx context.Context, blockID string) (*types.BeaconBlockHeader, error) {
	rv, err := c.getBlockHeader(ctx, blockID)
	if err != nil {
		c.logger.
			WithField("block", blockID).
			WithError(err).Errorf("GetBlockHeader failed")
	}

	return rv, err
}

// GetBlockHeader returns block header for given blockID
func (c *Client) getBlockHeader(ctx context.Context, blockID string) (*types.BeaconBlockHeader, error) {
	req, err := newGetBlockHeaderRequest(ctx, blockID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockHeader", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockHeader", resp, "Failure sending request")
	}

	result, err := inspectGetBlockHeaderResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockHeader", resp, "Invalid response")
	}

	return result, nil
}

func newGetBlockHeaderRequest(ctx context.Context, blockID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"blockID": autorest.Encode("path", blockID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/headers/{blockID}", pathParameters),
	).Prepare(newRequest(ctx))
}

type getBlockHeaderResponseMsg struct {
	Data *types.BeaconBlockHeader `json:"data"`
}

func inspectGetBlockHeaderResponse(resp *http.Response) (*types.BeaconBlockHeader, error) {
	msg := new(getBlockHeaderResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
