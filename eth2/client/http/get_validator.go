package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"

	"github.com/skillz-blockchain/go-utils/eth2/types"
)

// GetValidator returns validator specified by stateID and validatorID
func (c *Client) GetValidator(ctx context.Context, stateID, validatorID string) (*types.Validator, error) {
	req, err := newGetValidatorRequest(ctx, stateID, validatorID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidator", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidator", resp, "Failure sending request")
	}

	result, err := inspectGetValidatorResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetValidator", resp, "Invalid response")
	}

	return result, nil
}

func newGetValidatorRequest(ctx context.Context, stateID, validatorID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID":     autorest.Encode("path", stateID),
		"validatorID": autorest.Encode("path", validatorID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/validators/{validatorID}", pathParameters),
	).Prepare(newRequest(ctx))
}

type getValidatorResponseMsg struct {
	Data *types.Validator `json:"data"`
}

func inspectGetValidatorResponse(resp *http.Response) (*types.Validator, error) {
	msg := new(getValidatorResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
