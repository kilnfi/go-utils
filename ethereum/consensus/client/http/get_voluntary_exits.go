package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconphase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// GetVoluntaryExits returns voluntary exits known by the node but not necessarily incorporated into any block.
func (c *Client) GetVoluntaryExits(ctx context.Context) (beaconphase0.VoluntaryExits, error) {
	rv, err := c.getVoluntaryExits(ctx)
	if err != nil {
		c.logger.WithError(err).Errorf("GetVoluntaryExits failed")
	}

	return rv, err
}

func (c *Client) getVoluntaryExits(ctx context.Context) (beaconphase0.VoluntaryExits, error) {
	req, err := newGetVoluntaryExitsRequest(ctx)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetVoluntaryExits", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetVoluntaryExits", resp, "Failure sending request")
	}

	result, err := inspectGetVoluntaryExitsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetVoluntaryExits", resp, "Invalid response")
	}

	return result, nil
}

func newGetVoluntaryExitsRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("/eth/v1/beacon/pool/voluntary_exits"),
	).Prepare(newRequest(ctx))
}

type getVoluntaryExitsResponseMsg struct {
	Data beaconphase0.VoluntaryExits `json:"data"`
}

func inspectGetVoluntaryExitsResponse(resp *http.Response) (beaconphase0.VoluntaryExits, error) {
	msg := new(getVoluntaryExitsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
