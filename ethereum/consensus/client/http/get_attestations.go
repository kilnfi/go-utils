package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconphase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// GetAttestations returns attestations known by the node but not necessarily incorporated into any block.
func (c *Client) GetAttestations(ctx context.Context) (beaconphase0.Attestations, error) {
	req, err := newGetAttestationsRequest(ctx)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetAttestations", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetAttestations", resp, "Failure sending request")
	}

	result, err := inspectGetAttestationsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetAttestations", resp, "Invalid response")
	}

	return result, nil
}

func newGetAttestationsRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("eth/v1/beacon/attestations"),
	).Prepare(newRequest(ctx))
}

type getAttestationsResponseMsg struct {
	Data beaconphase0.Attestations `json:"data"`
}

func inspectGetAttestationsResponse(resp *http.Response) (beaconphase0.Attestations, error) {
	msg := new(getAttestationsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
