package eth2http

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"

	"github.com/skillz-blockchain/go-utils/eth2/types"
)

func WithBeaconErrorUnlessOK() autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			err := autorest.DecorateResponder(
				r,
				autorest.WithErrorUnlessOK(),
			).Respond(resp)
			if err == nil {
				return nil
			}

			if resp.StatusCode >= 400 {
				msg, beaconErr := inspectError(resp)
				if beaconErr == nil {
					err = autorest.NewErrorWithError(msg, "eth2http", "WitBeaconErrorUnlessOK", resp, "Failure with beacon node error")
				}
			}

			return err
		})
	}
}

func inspectError(resp *http.Response) (*types.Error, error) {
	msg := new(types.Error)
	err := autorest.Respond(
		resp,
		autorest.ByUnmarshallingJSON(msg),
		autorest.ByClosing(),
	)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
