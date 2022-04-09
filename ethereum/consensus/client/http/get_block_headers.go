package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"

	"github.com/skillz-blockchain/go-utils/ethereum/consensus/types"
)

// GetBlockHeaders return block headers
// Set slot and/or parentRoot to filter result (if nil no filter is applied)
func (c *Client) GetBlockHeaders(ctx context.Context, slot *beaconcommon.Slot, parentRoot *beaconcommon.Root) ([]*types.BeaconBlockHeader, error) {
	req, err := newGetBlockHeadersRequest(ctx, slot, parentRoot)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockHeaders", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockHeaders", resp, "Failure sending request")
	}

	result, err := inspectGetBlockHeadersResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockHeaders", resp, "Invalid response")
	}

	return result, nil
}

func newGetBlockHeadersRequest(ctx context.Context, slot *beaconcommon.Slot, parentRoot *beaconcommon.Root) (*http.Request, error) {
	queryParameters := map[string]interface{}{}
	if slot != nil {
		queryParameters["slot"] = slot.String()
	}

	if parentRoot != nil {
		queryParameters["parent_root"] = parentRoot.String()
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("eth/v1/beacon/headers"),
		autorest.WithQueryParameters(queryParameters),
	).Prepare(newRequest(ctx))
}

type getBlockHeadersResponseMsg struct {
	Data []*types.BeaconBlockHeader `json:"data"`
}

func inspectGetBlockHeadersResponse(resp *http.Response) ([]*types.BeaconBlockHeader, error) {
	msg := new(getBlockHeadersResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
