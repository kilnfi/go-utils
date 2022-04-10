package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

// GetBlockRoot returns hashTreeRoot of block
func (c *Client) GetBlockRoot(ctx context.Context, blockID string) (*beaconcommon.Root, error) {
	rv, err := c.getBlockRoot(ctx, blockID)
	if err != nil {
		c.logger.
			WithField("block", blockID).
			WithError(err).Errorf("GetBlockRoot failed")
	}

	return rv, err
}

func (c *Client) getBlockRoot(ctx context.Context, blockID string) (*beaconcommon.Root, error) {
	req, err := newGetBlockRootRequest(ctx, blockID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockRoot", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockRoot", resp, "Failure sending request")
	}

	result, err := InspectGetBlockRootResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockRoot", resp, "Invalid response")
	}

	return result, nil
}

func newGetBlockRootRequest(ctx context.Context, blockID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"blockID": autorest.Encode("path", blockID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/blocks/{blockID}/root", pathParameters),
	).Prepare(newRequest(ctx))
}

type getBlockRootResponseMsg struct {
	Data root `json:"data"`
}

func InspectGetBlockRootResponse(resp *http.Response) (*beaconcommon.Root, error) {
	msg := new(getBlockRootResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data.Root, nil
}
