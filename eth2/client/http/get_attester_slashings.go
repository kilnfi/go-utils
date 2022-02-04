package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconphase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// GetAttesterSlashings returns attester slashings known by the node but not necessarily incorporated into any block.
func (c *Client) GetAttesterSlashings(ctx context.Context) (beaconphase0.AttesterSlashings, error) {
	req, err := newGetAttesterSlashingsRequest(ctx)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetAttesterSlashings", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetAttesterSlashings", resp, "Failure sending request")
	}

	result, err := inspectGetAttesterSlashingsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetAttesterSlashings", resp, "Invalid response")
	}

	return result, nil
}

func newGetAttesterSlashingsRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("/eth/v1/beacon/pool/attester_slashings"),
	).Prepare(newRequest(ctx))
}

type getAttesterSlashingsResponseMsg struct {
	Data beaconphase0.AttesterSlashings `json:"data"`
}

func inspectGetAttesterSlashingsResponse(resp *http.Response) (beaconphase0.AttesterSlashings, error) {
	msg := new(getAttesterSlashingsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
