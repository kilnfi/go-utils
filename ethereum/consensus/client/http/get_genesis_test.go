package eth2http

import (
	"context"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skillz-blockchain/go-utils/ethereum/consensus/types"
	httptestutils "github.com/skillz-blockchain/go-utils/http/testutils"
)

func TestGetGenesis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	t.Run("StatusOK", func(t *testing.T) { testGetGenesisStatusOK(t, c, mockCli) })
	t.Run("Status400", func(t *testing.T) { testGetGenesisStatus400(t, c, mockCli) })
}

func testGetGenesisStatusOK(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Get("/eth/v1/beacon/genesis").
		Reply(200).
		JSON([]byte(`{"data":{"genesis_time":"1606824023","genesis_validators_root":"0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95","genesis_fork_version":"0x00000000"}}`))

	mockCli.EXPECT().Gock(req)

	genesis, err := c.GetGenesis(context.Background())

	require.NoError(t, err)
	assert.Equal(
		t,
		&types.Genesis{
			GenesisTime:           beaconcommon.Timestamp(1606824023),
			GenesisValidatorsRoot: beaconcommon.Root(gethcommon.HexToHash("0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95")),
		},
		genesis,
	)
}

func testGetGenesisStatus400(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Get("/eth/v1/beacon/genesis").
		Reply(400)

	mockCli.EXPECT().Gock(req)

	_, err := c.GetGenesis(context.Background())

	require.Error(t, err)
}
