package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

// GetStateRoot calculates HashTreeRoot of the state for the given stateID
func (c *Client) GetStateRoot(ctx context.Context, stateID string) (*beaconcommon.Root, error) {
	req, err := newGetStateRootRequest(ctx, stateID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateRoot", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateRoot", resp, "Failure sending request")
	}

	result, err := inspectGetStateRootResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetStateRoot", resp, "Invalid response")
	}

	return result, nil
}

func newGetStateRootRequest(ctx context.Context, stateID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/root", pathParameters),
	).Prepare(newRequest(ctx))
}

type root struct {
	Root *beaconcommon.Root `json:"root"`
}

type getStateRootResponseMsg struct {
	Data root `json:"data"`
}

func inspectGetStateRootResponse(resp *http.Response) (*beaconcommon.Root, error) {
	msg := new(getStateRootResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data.Root, nil
}
