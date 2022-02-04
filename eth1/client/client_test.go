package eth1

import (
	"context"
	"math/big"
	"testing"

	geth "github.com/ethereum/go-ethereum"
	gethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httptestutils "github.com/skillz-blockchain/go-utils/http/testutils"
	jsonrpchttp "github.com/skillz-blockchain/go-utils/jsonrpc/http"
)

func TestClientImplementsGetBindingInterface(t *testing.T) {
	client := new(Client)
	assert.Implements(t, (*gethbind.ContractCaller)(nil), client)
	assert.Implements(t, (*gethbind.ContractTransactor)(nil), client)
	assert.Implements(t, (*gethbind.ContractFilterer)(nil), client)
}

func TestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := New(jsonrpchttp.NewClient(mockCli))

	t.Run("BlockNumber", func(t *testing.T) { testBlockNumber(t, c, mockCli) })
	t.Run("HeaderByNumber", func(t *testing.T) { testHeaderByNumber(t, c, mockCli) })
	t.Run("CallContract", func(t *testing.T) { testCallContract(t, c, mockCli) })
	t.Run("NonceAt", func(t *testing.T) { testNonceAt(t, c, mockCli) })
	t.Run("PendingNonceAt", func(t *testing.T) { testPendingNonceAt(t, c, mockCli) })
	t.Run("SuggestGasPrice", func(t *testing.T) { testSuggestGasPrice(t, c, mockCli) })
	t.Run("SuggestGasTipCap", func(t *testing.T) { testSuggestGasTipCap(t, c, mockCli) })
	t.Run("EstimateGas", func(t *testing.T) { testEstimateGas(t, c, mockCli) })
	t.Run("SendTransaction", func(t *testing.T) { testSendTransaction(t, c, mockCli) })
}

func testBlockNumber(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_blockNumber","params":null,"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","result":"0x20","id":0}`))

	mockCli.EXPECT().Gock(req)

	blockNumber, err := c.BlockNumber(context.Background())

	require.NoError(t, err)
	assert.Equal(t, uint64(32), blockNumber)
}

func testHeaderByNumber(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_getBlockByNumber","params":["0xd6e166",false],"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","id":1,"result":{"baseFeePerGas":"0x1c30017ca8","difficulty":"0x2d754c4a5c3f14","extraData":"0x6e616e6f706f6f6c2e6f7267","gasLimit":"0x1c9c364","gasUsed":"0x1c985bc","hash":"0x0fb6d5609c9edab75bf587ea7449e6e6940d6e3df1992a1bd96ca8b74ffd16fc","logsBloom":"0x1fbb5f53e8e63cedffe45fd8bf1217fdee15d39bbebf275136afb8ffb99fdd9b92556ffb2ceeb1345a3bf1dd730ebfc6bf4c814119e6faaef2f9fa9b50ffe8fd838eb2bed773592efb0ffc7efd142fe37fe65117f5f4f7bb2f037671a4ff52d443a7044a1be25ec1fb1b13a9aabf6afdd278f4bf4abda64e3293cb9480f97d11c9558ded275cdf8ed5ef7f43398e9fb5fe4e2e0d79257cecebf95bd36e99a8f7bbdab5323febe6baceb1dfdda71cbe21dfbcc6a3feee6702fd85a6bd3ee9f8dc757ca4bacdf3a47ef119c3d95feb5d2f65acffdb9effa17ebb5fdb1b3afe64dfd8fcf3bfa8787f882e660d33cfe7fb9220ef6226efd5dffafcc7daa3b6967faf","miner":"0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","mixHash":"0x274264e3a69256c43beb4632b6bf8ac2de6534dd6c4fb09dad1a0541eb8ed356","nonce":"0x2fdaedd11fd5a2ea","number":"0xd6e166","parentHash":"0x6019a4b3e4e3ba7b7b43d28d68492f99226b86e7dff0c607a16ef4d16a617503","receiptsRoot":"0x081119bc627ccedade0b6321984146672ad1a15b0769b08f7a91ea22474c7bd9","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x43a1d","stateRoot":"0x4a4e5f11b8e837adb24fb764ab93f33ed21efa279df4fe59b5bed3c3885e9fae","timestamp":"0x61f179e3","totalDifficulty":"0x873cd0f1a366947ae8d","transactions":[],"transactionsRoot":"0x5cb8acbd8a0d2f3c489e47d8267c86a718203da8a5a34f0511918c13cbb14c1b","uncles":[]},"id":0}`))

	mockCli.EXPECT().Gock(req)

	header, err := c.HeaderByNumber(context.Background(), big.NewInt(14082406))

	require.NoError(t, err)
	assert.Equal(
		t,
		&gethtypes.Header{
			ParentHash:  gethcommon.HexToHash("0x6019a4b3e4e3ba7b7b43d28d68492f99226b86e7dff0c607a16ef4d16a617503"),
			UncleHash:   gethcommon.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
			Coinbase:    gethcommon.HexToAddress("0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5"),
			Root:        gethcommon.HexToHash("0x4a4e5f11b8e837adb24fb764ab93f33ed21efa279df4fe59b5bed3c3885e9fae"),
			TxHash:      gethcommon.HexToHash("0x5cb8acbd8a0d2f3c489e47d8267c86a718203da8a5a34f0511918c13cbb14c1b"),
			ReceiptHash: gethcommon.HexToHash("0x081119bc627ccedade0b6321984146672ad1a15b0769b08f7a91ea22474c7bd9"),
			Bloom:       gethtypes.BytesToBloom(gethcommon.FromHex("0x1fbb5f53e8e63cedffe45fd8bf1217fdee15d39bbebf275136afb8ffb99fdd9b92556ffb2ceeb1345a3bf1dd730ebfc6bf4c814119e6faaef2f9fa9b50ffe8fd838eb2bed773592efb0ffc7efd142fe37fe65117f5f4f7bb2f037671a4ff52d443a7044a1be25ec1fb1b13a9aabf6afdd278f4bf4abda64e3293cb9480f97d11c9558ded275cdf8ed5ef7f43398e9fb5fe4e2e0d79257cecebf95bd36e99a8f7bbdab5323febe6baceb1dfdda71cbe21dfbcc6a3feee6702fd85a6bd3ee9f8dc757ca4bacdf3a47ef119c3d95feb5d2f65acffdb9effa17ebb5fdb1b3afe64dfd8fcf3bfa8787f882e660d33cfe7fb9220ef6226efd5dffafcc7daa3b6967faf")),
			Difficulty:  big.NewInt(12795344477503252),
			Number:      big.NewInt(14082406),
			GasLimit:    29999972,
			GasUsed:     29984188,
			Time:        uint64(1643215331),
			Extra:       gethcommon.FromHex("0x6e616e6f706f6f6c2e6f7267"),
			MixDigest:   gethcommon.HexToHash("0x274264e3a69256c43beb4632b6bf8ac2de6534dd6c4fb09dad1a0541eb8ed356"),
			Nonce:       gethtypes.EncodeNonce(3448329947143578346),
			BaseFee:     big.NewInt(121064488104),
		},
		header,
	)
}

func testCallContract(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_call","params":[{"data":"0x0123456789","from":"0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","to":null},"0xd6e166"],"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","result":"0xabcdef","id":0}`))

	mockCli.EXPECT().Gock(req)

	res, err := c.CallContract(
		context.Background(),
		geth.CallMsg{
			From: gethcommon.HexToAddress("0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5"),
			Data: gethcommon.FromHex("0x0123456789"),
		},
		big.NewInt(14082406),
	)

	require.NoError(t, err)
	assert.Equal(t, gethcommon.FromHex("0xabcdef"), res)
}

func testNonceAt(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_getTransactionCount","params":["0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","0xd6e6f3"],"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x11189c7"}`))

	mockCli.EXPECT().Gock(req)

	nonce, err := c.NonceAt(
		context.Background(),
		gethcommon.HexToAddress("0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5"),
		big.NewInt(14083827),
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(17926599), nonce)
}

func testPendingNonceAt(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_getTransactionCount","params":["0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","pending"],"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x11189c7"}`))

	mockCli.EXPECT().Gock(req)

	nonce, err := c.PendingNonceAt(
		context.Background(),
		gethcommon.HexToAddress("0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5"),
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(17926599), nonce)
}

func testSuggestGasPrice(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_gasPrice","params":null,"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x2fbbd1aa9b"}`))

	mockCli.EXPECT().Gock(req)

	p, err := c.SuggestGasPrice(context.Background())

	require.NoError(t, err)
	assert.Equal(t, big.NewInt(205014543003), p)
}

func testSuggestGasTipCap(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_maxPriorityFeePerGas","params":null,"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x2fbbd1aa9b"}`))

	mockCli.EXPECT().Gock(req)

	p, err := c.SuggestGasTipCap(context.Background())

	require.NoError(t, err)
	assert.Equal(t, big.NewInt(205014543003), p)
}

func testEstimateGas(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_estimateGas","params":[{"data":"0x0123456789","from":"0x52bc44d5378309ee2abf1539bf71de1b7d7be3b5","to":null}],"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","result":"0xabcdef","id":0}`))

	mockCli.EXPECT().Gock(req)

	gas, err := c.EstimateGas(
		context.Background(),
		geth.CallMsg{
			From: gethcommon.HexToAddress("0x52bc44d5378309EE2abF1539BF71dE1b7d7bE3b5"),
			Data: gethcommon.FromHex("0x0123456789"),
		},
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(11259375), gas)
}

func testSendTransaction(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Post("/").
		JSON([]byte(`{"jsonrpc":"","method":"eth_sendRawTransaction","params":["0xf86d8202b38477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ba0699ff162205967ccbabae13e07cdd4284258d46ec1051a70a51be51ec2bc69f3a04e6944d508244ea54a62ebf9a72683eeadacb73ad7c373ee542f1998147b220e"],"id":null}`)).
		Reply(200).
		JSON([]byte(`{"jsonrpc":"2.0","result":"0x679bdd54941acaebcf592035101606b56087048ebb7ea12a02df4a6be426f8dd","id":0}`))

	mockCli.EXPECT().Gock(req)

	tx := &gethtypes.Transaction{}
	_ = tx.UnmarshalBinary(gethcommon.FromHex("0xf86d8202b38477359400825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d880de0b6b3a7640000802ba0699ff162205967ccbabae13e07cdd4284258d46ec1051a70a51be51ec2bc69f3a04e6944d508244ea54a62ebf9a72683eeadacb73ad7c373ee542f1998147b220e"))

	err := c.SendTransaction(
		context.Background(),
		tx,
	)

	require.NoError(t, err)
}
