package eth2http

import (
	"context"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"

	"github.com/skillz-blockchain/go-utils/eth2/types"
)

// GetValidators returns list of validators
// Set validatorsIDs and/or statuses to filter result (if empty no filter is applied)
func (c *Client) GetValidators(ctx context.Context, stateID string, validatorIDs, statuses []string) ([]*types.Validator, error) {
	req, err := newGetValidatorsRequest(ctx, stateID, validatorIDs, statuses)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidators", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidators", resp, "Failure sending request")
	}

	result, err := inspectGetValidatorsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidators", resp, "Invalid response")
	}

	return result, nil
}

func newGetValidatorsRequest(ctx context.Context, stateID string, validatorIDs, statuses []string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	queryParameters := map[string]interface{}{}
	if len(validatorIDs) != 0 {
		queryParameters["id"] = strings.Join(validatorIDs, ",")
	}

	if len(statuses) != 0 {
		queryParameters["status"] = strings.Join(statuses, ",")
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/validators", pathParameters),
		autorest.WithQueryParameters(queryParameters),
	).Prepare(newRequest(ctx))
}

type getValidatorsResponseMsg struct {
	Data []*types.Validator `json:"data"`
}

func inspectGetValidatorsResponse(resp *http.Response) ([]*types.Validator, error) {
	msg := new(getValidatorsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
