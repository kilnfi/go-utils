package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"

	"github.com/skillz-blockchain/go-utils/ethereum/consensus/types"
)

// GetStateFinalityCheckpoints returns finality checkpoints for state with given stateID
// In case finality is not yet achieved returns epoch 0 and ZERO_HASH as root.
func (c *Client) GetStateFinalityCheckpoints(ctx context.Context, stateID string) (*types.StateFinalityCheckpoints, error) {
	rv, err := c.getStateFinalityCheckpoints(ctx, stateID)
	if err != nil {
		c.logger.
			WithField("state", stateID).
			WithError(err).Errorf("GetStateFinalityCheckpoints failed")
	}

	return rv, err
}

func (c *Client) getStateFinalityCheckpoints(ctx context.Context, stateID string) (*types.StateFinalityCheckpoints, error) {
	req, err := newGetStateFinalityCheckpointsRequest(ctx, stateID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateFinalityCheckpoints", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateFinalityCheckpoints", resp, "Failure sending request")
	}

	result, err := inspectGetStateFinalityCheckpointsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateFinalityCheckpoints", resp, "Invalid response")
	}

	return result, nil
}

func newGetStateFinalityCheckpointsRequest(ctx context.Context, stateID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}}/finality_checkpoints", pathParameters),
	).Prepare(newRequest(ctx))
}

type getStateFinalityCheckpointsResponseMsg struct {
	Data *types.StateFinalityCheckpoints `json:"data"`
}

func inspectGetStateFinalityCheckpointsResponse(resp *http.Response) (*types.StateFinalityCheckpoints, error) {
	msg := new(getStateFinalityCheckpointsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
