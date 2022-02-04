package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"

	"github.com/skillz-blockchain/go-utils/eth2/types"
)

// GetSyncCommittees returns the sync committees for given stateID
// Set epoch to filter result (if nil no filter is applied)
func (c *Client) GetSyncCommittees(ctx context.Context, stateID string, epoch *beaconcommon.Epoch) (*types.SyncCommittees, error) {
	req, err := newGetSyncCommitteesRequest(ctx, stateID, epoch)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetSyncCommittees", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetSyncCommittees", resp, "Failure sending request")
	}

	result, err := inspectGetSyncCommitteesResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetSyncCommittees", resp, "Invalid response")
	}

	return result, nil
}

func newGetSyncCommitteesRequest(ctx context.Context, stateID string, epoch *beaconcommon.Epoch) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	queryParameters := map[string]interface{}{}
	if epoch != nil {
		queryParameters["epoch"] = epoch.String()
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/sync_committees", pathParameters),
		autorest.WithQueryParameters(queryParameters),
	).Prepare(newRequest(ctx))
}

type getSyncCommitteesResponseMsg struct {
	Data *types.SyncCommittees `json:"data"`
}

func inspectGetSyncCommitteesResponse(resp *http.Response) (*types.SyncCommittees, error) {
	msg := new(getSyncCommitteesResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
