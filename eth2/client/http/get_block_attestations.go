package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconphase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// GetBlockAttestations returns attestations included in requested block with given blockID
func (c *Client) GetBlockAttestations(ctx context.Context, blockID string) (beaconphase0.Attestations, error) {
	req, err := newGetBlockAttestationsRequest(ctx, blockID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockAttestations", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockAttestations", resp, "Failure sending request")
	}

	result, err := inspectGetBlockAttestationsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockAttestations", resp, "Invalid response")
	}

	return result, nil
}

func newGetBlockAttestationsRequest(ctx context.Context, blockID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"blockID": autorest.Encode("path", blockID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/blocks/{blockID}/attestations", pathParameters),
	).Prepare(newRequest(ctx))
}

type getBlockAttestationsResponseMsg struct {
	Data beaconphase0.Attestations `json:"data"`
}

func inspectGetBlockAttestationsResponse(resp *http.Response) (beaconphase0.Attestations, error) {
	msg := new(getBlockAttestationsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
