package eth2http

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skillz-blockchain/go-utils/ethereum/consensus/types"
	httptestutils "github.com/skillz-blockchain/go-utils/net/http/testutils"
)

func TestGetValidators(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	t.Run("StatusOK", func(t *testing.T) { testGetValidatorsStatusOK(t, c, mockCli) })
}

func testGetValidatorsStatusOK(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Get("/eth/v1/beacon/states/test-state/validators").
		MatchParams(map[string]string{
			"status": "sA,sB",
			"id":     "vA,vB,vC",
		}).
		Reply(200).
		JSON([]byte(`{"data":[]}`))

	mockCli.EXPECT().Gock(req)

	vals, err := c.GetValidators(context.Background(), "test-state", []string{"vA", "vB", "vC"}, []string{"sA", "sB"})
	require.NoError(t, err)
	assert.Equal(
		t,
		[]*types.Validator{},
		vals,
	)
}
