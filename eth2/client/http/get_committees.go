package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"

	"github.com/skillz-blockchain/go-utils/eth2/types"
)

// GetCommittees returns the committees for the given state.
// Set epoch and/or index and/or slot to filter result (if nil no filter is applied)
func (c *Client) GetCommittees(ctx context.Context, stateID string, epoch *beaconcommon.Epoch, index *beaconcommon.CommitteeIndex, slot *beaconcommon.Slot) ([]*types.Committee, error) {
	req, err := newGetCommitteesRequest(ctx, stateID, epoch, index, slot)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetCommittees", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetCommittees", resp, "Failure sending request")
	}

	result, err := inspectGetCommitteesResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetCommittees", resp, "Invalid response")
	}

	return result, nil
}

func newGetCommitteesRequest(ctx context.Context, stateID string, epoch *beaconcommon.Epoch, index *beaconcommon.CommitteeIndex, slot *beaconcommon.Slot) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	queryParameters := map[string]interface{}{}
	if epoch != nil {
		queryParameters["epoch"] = epoch.String()
	}

	if index != nil {
		queryParameters["index"] = index.String()
	}

	if slot != nil {
		queryParameters["slot"] = slot.String()
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/committees", pathParameters),
		autorest.WithQueryParameters(queryParameters),
	).Prepare(newRequest(ctx))
}

type getCommitteesResponseMsg struct {
	Data []*types.Committee `json:"data"`
}

func inspectGetCommitteesResponse(resp *http.Response) ([]*types.Committee, error) {
	msg := new(getCommitteesResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
