package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconphase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// GetProposerSlashings returns proposer slashings known by the node but not necessarily incorporated into any block.
func (c *Client) GetProposerSlashings(ctx context.Context) (beaconphase0.ProposerSlashings, error) {
	rv, err := c.getProposerSlashings(ctx)
	if err != nil {
		c.logger.WithError(err).Errorf("GetProposerSlashings failed")
	}

	return rv, err
}

func (c *Client) getProposerSlashings(ctx context.Context) (beaconphase0.ProposerSlashings, error) {
	req, err := newGetProposerSlashingsRequest(ctx)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetProposerSlashings", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetProposerSlashings", resp, "Failure sending request")
	}

	result, err := inspectGetProposerSlashingsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetProposerSlashings", resp, "Invalid response")
	}

	return result, nil
}

func newGetProposerSlashingsRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("/eth/v1/beacon/pool/proposer_slashings"),
	).Prepare(newRequest(ctx))
}

type getProposerSlashingsResponseMsg struct {
	Data beaconphase0.ProposerSlashings `json:"data"`
}

func inspectGetProposerSlashingsResponse(resp *http.Response) (beaconphase0.ProposerSlashings, error) {
	msg := new(getProposerSlashingsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
