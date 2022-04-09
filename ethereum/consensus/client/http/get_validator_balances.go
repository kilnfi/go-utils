package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"

	"github.com/skillz-blockchain/go-utils/ethereum/consensus/types"
)

// GetValidatorBalances returns list of validator balances.
// Set validatorsIDs to filter validator result (if empty no filter is applied)
func (c *Client) GetValidatorBalances(ctx context.Context, stateID string, validatorIDs []string) ([]*types.ValidatorBalance, error) {
	req, err := newGetValidatorBalancesRequest(ctx, stateID, validatorIDs)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidatorBalances", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidatorBalances", resp, "Failure sending request")
	}

	result, err := inspectGetValidatorBalancesResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidatorBalances", resp, "Invalid response")
	}

	return result, nil
}

func newGetValidatorBalancesRequest(ctx context.Context, stateID string, validatorIDs []string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	queryParameters := map[string]interface{}{}
	if len(validatorIDs) != 0 {
		queryParameters["validator_id"] = validatorIDs
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/ValidatorBalances", pathParameters),
		autorest.WithQueryParameters(queryParameters),
	).Prepare(newRequest(ctx))
}

type getValidatorBalancesResponseMsg struct {
	Data []*types.ValidatorBalance `json:"data"`
}

func inspectGetValidatorBalancesResponse(resp *http.Response) ([]*types.ValidatorBalance, error) {
	msg := new(getValidatorBalancesResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
