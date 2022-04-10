package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

// GetStateFork returns Fork object for state with given stateID
func (c *Client) GetStateFork(ctx context.Context, stateID string) (*beaconcommon.Fork, error) {
	rv, err := c.getStateFork(ctx, stateID)
	if err != nil {
		c.logger.
			WithField("state", stateID).
			WithError(err).Errorf("GetStateFork failed")
	}

	return rv, err
}

func (c *Client) getStateFork(ctx context.Context, stateID string) (*beaconcommon.Fork, error) {
	req, err := newGetStateForkRequest(ctx, stateID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateFork", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateFork", resp, "Failure sending request")
	}

	result, err := inspectGetStateForkResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateFork", resp, "Invalid response")
	}

	return result, nil
}

func newGetStateForkRequest(ctx context.Context, stateID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/fork", pathParameters),
	).Prepare(newRequest(ctx))
}

type getStateForkResponseMsg struct {
	Data *beaconcommon.Fork `json:"data"`
}

func inspectGetStateForkResponse(resp *http.Response) (*beaconcommon.Fork, error) {
	msg := new(getStateForkResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
