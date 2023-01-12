package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	geth "github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"

	"github.com/kilnfi/go-utils/common/interfaces"
	"github.com/kilnfi/go-utils/ethereum/execution/client"
	"github.com/kilnfi/go-utils/ethereum/execution/types"
	"github.com/kilnfi/go-utils/net/jsonrpc"
	jsonrpchttp "github.com/kilnfi/go-utils/net/jsonrpc/http"
)

// Ensure Client interface is fully implemented
var _ client.Client = (*Client)(nil)

// Client provides methods to interface with a JSON-RPC Ethereum 1.0 node
type Client struct {
	client jsonrpc.Client

	chainID *big.Int
	mu      sync.Mutex
}

// New creates a new client
func NewFromClient(cli jsonrpc.Client) *Client {
	return &Client{
		client: cli,
	}
}

// NewFromAddress creates a new client connecting to an Ethereum node at addr
func New(cfg *jsonrpchttp.Config) (*Client, error) {
	jsonrpcc, err := jsonrpchttp.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return NewFromClient(jsonrpc.WithIncrementalID()(jsonrpc.WithVersion("2.0")(jsonrpcc))), nil
}

func (c *Client) Logger() logrus.FieldLogger {
	if loggable, ok := c.client.(interfaces.Loggable); ok {
		return loggable.Logger()
	}
	return nil
}

func (c *Client) SetLogger(logger logrus.FieldLogger) {
	if loggable, ok := c.client.(interfaces.Loggable); ok {
		loggable.SetLogger(logger)
	}
}

func (c *Client) call(ctx context.Context, res interface{}, method string, params ...interface{}) error {
	return c.client.Call(
		ctx,
		&jsonrpc.Request{
			Method: method,
			Params: params,
		},
		res,
	)
}

// ChainID retrieves the current chain ID
func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.chainID != nil {
		return c.chainID, nil
	}
	res := new(gethhexutil.Big)
	err := c.call(ctx, res, "eth_chainId")
	if err != nil {
		return nil, err
	}
	c.chainID = (*big.Int)(res)
	return c.chainID, nil
}

// BlockNumber returns current chain head number
func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	res := new(gethhexutil.Uint64)
	err := c.call(ctx, res, "eth_blockNumber")
	if err != nil {
		return 0, err
	}

	return uint64(*res), nil
}

func (c *Client) BalanceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	res := new(gethhexutil.Big)
	err := c.call(ctx, res, "eth_getBalance", account, types.ToBlockNumArg(blockNumber))
	return (*big.Int)(res), err
}

// BlockByHash returns the given full block.
//
// Note fetch of uncles blocks is not implemented yet.
func (c *Client) BlockByHash(ctx context.Context, hash gethcommon.Hash) (*gethtypes.Block, error) {
	return c.getBlock(ctx, "eth_getBlockByHash", hash, true)
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note fetch of uncles blocks is not implemented yet.
func (c *Client) BlockByNumber(ctx context.Context, blockNumber *big.Int) (*gethtypes.Block, error) {
	return c.getBlock(ctx, "eth_getBlockByNumber", types.ToBlockNumArg(blockNumber), true)
}

// HeaderByNumber returns header of a given block hash
func (c *Client) HeaderByHash(ctx context.Context, hash gethcommon.Hash) (*gethtypes.Header, error) {
	res := new(gethtypes.Header)
	err := c.call(ctx, res, "eth_getBlockByHash", hash, false)
	if err == nil && res == nil {
		err = geth.NotFound
	}
	return res, err
}

// HeaderByNumber returns header of a given block number
func (c *Client) HeaderByNumber(ctx context.Context, blockNumber *big.Int) (*gethtypes.Header, error) {
	res := new(gethtypes.Header)
	err := c.call(ctx, res, "eth_getBlockByNumber", types.ToBlockNumArg(blockNumber), false)
	if err == nil && res == nil {
		err = geth.NotFound
	}

	return res, err
}

// CallContract executes contract call
// The block number can be nil, in which case call is executed at the latest block.
//
//nolint:gocritic
func (c *Client) CallContract(ctx context.Context, msg geth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	res := new(gethhexutil.Bytes)
	err := c.call(ctx, res, "eth_call", toCallArg(&msg), types.ToBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	return []byte(*res), nil
}

// CallContractAtHash is almost the same as CallContract except that it selects
// the block by block hash instead of block height.
//
//nolint:gocritic
func (c *Client) CallContractAtHash(ctx context.Context, msg geth.CallMsg, blockHash gethcommon.Hash) ([]byte, error) {
	var res gethhexutil.Bytes
	err := c.call(ctx, res, "eth_call", toCallArg(&msg), gethrpc.BlockNumberOrHashWithHash(blockHash, false))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest block.
func (c *Client) CodeAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	res := new(gethhexutil.Bytes)
	err := c.call(ctx, res, "eth_getCode", account, types.ToBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	return []byte(*res), nil
}

// PendingCodeAt returns the contract code of the given account on pending state
func (c *Client) PendingCodeAt(ctx context.Context, account gethcommon.Address) ([]byte, error) {
	return c.CodeAt(ctx, account, big.NewInt(-1))
}

// NonceAt returns the next nonce for the given account.
// The block number can be nil, in which case the code is taken from the latest block.
func (c *Client) NonceAt(ctx context.Context, account gethcommon.Address, blockNumber *big.Int) (uint64, error) {
	res := new(gethhexutil.Uint64)
	err := c.call(ctx, res, "eth_getTransactionCount", account, types.ToBlockNumArg(blockNumber))
	if err != nil {
		return 0, err
	}

	return uint64(*res), nil
}

// PendingNonceAt returns the next nonce for the given account considering pending transaction.
func (c *Client) PendingNonceAt(ctx context.Context, account gethcommon.Address) (uint64, error) {
	return c.NonceAt(ctx, account, big.NewInt(-1))
}

// SuggestGasPrice returns gas price for a transaction to be included in a miner block in a timely
// manner considering current network activity
func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	res := new(gethhexutil.Big)
	err := c.call(ctx, res, "eth_gasPrice")
	if err != nil {
		return nil, err
	}

	return (*big.Int)(res), nil
}

// SuggestGasPrice returns a gas tip cap after EIP-1559 for a transaction to be included in a miner block in a timely
// manner considering current network activity
func (c *Client) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	res := new(gethhexutil.Big)
	err := c.call(ctx, res, "eth_maxPriorityFeePerGas")
	if err != nil {
		return nil, err
	}

	return (*big.Int)(res), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the chain.
//
//nolint:gocritic
func (c *Client) EstimateGas(ctx context.Context, msg geth.CallMsg) (uint64, error) {
	res := new(gethhexutil.Uint64)
	err := c.call(ctx, res, "eth_estimateGas", toCallArg(&msg))
	if err != nil {
		return 0, err
	}
	return uint64(*res), nil
}

// SendTransaction injects a signed transaction into the pending pool for execution.
func (c *Client) SendTransaction(ctx context.Context, tx *gethtypes.Transaction) error {
	data, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	return c.call(ctx, nil, "eth_sendRawTransaction", gethhexutil.Encode(data))
}

// FilterLogs executes a filter query.
func (c *Client) FilterLogs(ctx context.Context, q geth.FilterQuery) ([]gethtypes.Log, error) {
	var res []gethtypes.Log
	arg, err := toFilterArg(q)
	if err != nil {
		return nil, err
	}

	err = c.call(ctx, res, "eth_getLogs", arg)

	return res, err
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (c *Client) SubscribeFilterLogs(ctx context.Context, _ geth.FilterQuery, _ chan<- gethtypes.Log) (geth.Subscription, error) {
	return nil, fmt.Errorf("not implemented")
}

type feeHistoryResultMarshaling struct {
	OldestBlock  *gethhexutil.Big     `json:"oldestBlock"`
	Reward       [][]*gethhexutil.Big `json:"reward,omitempty"`
	BaseFee      []*gethhexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64            `json:"gasUsedRatio"`
}

// FeeHistory retrieves the fee market history.
func (c *Client) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*geth.FeeHistory, error) {
	var res feeHistoryResultMarshaling
	if err := c.call(ctx, &res, "eth_feeHistory", gethhexutil.Uint(blockCount), types.ToBlockNumArg(lastBlock), rewardPercentiles); err != nil {
		return nil, err
	}
	reward := make([][]*big.Int, len(res.Reward))
	for i, r := range res.Reward {
		reward[i] = make([]*big.Int, len(r))
		for j, r := range r {
			reward[i][j] = (*big.Int)(r)
		}
	}
	baseFee := make([]*big.Int, len(res.BaseFee))
	for i, b := range res.BaseFee {
		baseFee[i] = (*big.Int)(b)
	}
	return &geth.FeeHistory{
		OldestBlock:  (*big.Int)(res.OldestBlock),
		Reward:       reward,
		BaseFee:      baseFee,
		GasUsedRatio: res.GasUsedRatio,
	}, nil
}

// NetworkID returns the network ID (also known as the chain ID) for this chain.
func (c *Client) NetworkID(ctx context.Context) (*big.Int, error) {
	version := new(big.Int)
	var ver string
	if err := c.call(ctx, &ver, "net_version"); err != nil {
		return nil, err
	}
	if _, ok := version.SetString(ver, 10); !ok {
		return nil, fmt.Errorf("invalid net_version result %q", ver)
	}
	return version, nil
}

// PeerCount returns the number of p2p peers as reported by the net_peerCount method.
func (c *Client) PeerCount(ctx context.Context) (uint64, error) {
	var result gethhexutil.Uint64
	err := c.call(ctx, &result, "net_peerCount")
	return uint64(result), err
}

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (c *Client) PendingBalanceAt(ctx context.Context, account gethcommon.Address) (*big.Int, error) {
	var result gethhexutil.Big
	err := c.call(ctx, &result, "eth_getBalance", account, "pending")
	return (*big.Int)(&result), err
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
//
//nolint:gocritic
func (c *Client) PendingCallContract(ctx context.Context, msg geth.CallMsg) ([]byte, error) {
	var hex gethhexutil.Bytes
	err := c.call(ctx, &hex, "eth_call", toCallArg(&msg), "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (c *Client) PendingStorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash) ([]byte, error) {
	var result gethhexutil.Bytes
	err := c.call(ctx, &result, "eth_getStorageAt", account, key, "pending")
	return result, err
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (c *Client) PendingTransactionCount(ctx context.Context) (uint, error) {
	var num gethhexutil.Uint
	err := c.call(ctx, &num, "eth_getBlockTransactionCountByNumber", "pending")
	return uint(num), err
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (c *Client) StorageAt(ctx context.Context, account gethcommon.Address, key gethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	var result gethhexutil.Bytes
	err := c.call(ctx, &result, "eth_getStorageAt", account, key, types.ToBlockNumArg(blockNumber))
	return result, err
}

func (c *Client) SubscribeNewHead(ctx context.Context, ch chan<- *gethtypes.Header) (geth.Subscription, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) SyncProgress(ctx context.Context) (*geth.SyncProgress, error) {
	return nil, fmt.Errorf("not implemented")
}

// TransactionByHash returns the transaction with the given hash.
func (c *Client) TransactionByHash(ctx context.Context, hash gethcommon.Hash) (tx *gethtypes.Transaction, isPending bool, err error) {
	var res *types.RPCTransaction
	err = c.call(ctx, &res, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, false, err
	} else if res == nil {
		return nil, false, geth.NotFound
	} else if _, r, _ := res.Tx.RawSignatureValues(); r == nil {
		return nil, false, fmt.Errorf("server returned transaction without signature")
	}
	if res.From != nil && res.BlockHash != nil {
		setSenderFromServer(res.Tx, *res.From, *res.BlockHash)
	}
	return res.Tx, res.BlockNumber == nil, nil
}

// TransactionCount returns the total number of transactions in the given block.
func (c *Client) TransactionCount(ctx context.Context, blockHash gethcommon.Hash) (uint, error) {
	var num gethhexutil.Uint
	err := c.call(ctx, &num, "eth_getBlockTransactionCountByHash", blockHash)
	return uint(num), err
}

// TransactionInBlock returns a single transaction at index in the given block.
func (c *Client) TransactionInBlock(ctx context.Context, blockHash gethcommon.Hash, index uint) (*gethtypes.Transaction, error) {
	var res *types.RPCTransaction
	err := c.call(ctx, &res, "eth_getTransactionByBlockHashAndIndex", blockHash, gethhexutil.Uint64(index))
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, geth.NotFound
	} else if _, r, _ := res.Tx.RawSignatureValues(); r == nil {
		return nil, fmt.Errorf("server returned transaction without signature")
	}
	if res.From != nil && res.BlockHash != nil {
		setSenderFromServer(res.Tx, *res.From, *res.BlockHash)
	}
	return res.Tx, err
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (c *Client) TransactionReceipt(ctx context.Context, txHash gethcommon.Hash) (*gethtypes.Receipt, error) {
	var r *gethtypes.Receipt
	err := c.call(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, geth.NotFound
		}
	}
	return r, err
}

// TransactionSender returns the sender address of the given transaction. The transaction
// must be known to the remote node and included in the blockchain at the given block and
// index. The sender is the one derived by the protocol at the time of inclusion.
//
// There is a fast-path for transactions retrieved by TransactionByHash and
// TransactionInBlock. Getting their sender address can be done without an RPC interaction.
func (c *Client) TransactionSender(ctx context.Context, tx *gethtypes.Transaction, block gethcommon.Hash, index uint) (gethcommon.Address, error) {
	// Try to load the address from the cache.
	sender, err := gethtypes.Sender(&senderFromServer{blockhash: block}, tx)
	if err == nil {
		return sender, nil
	}

	// It was not found in cache, ask the server.
	var meta struct {
		Hash gethcommon.Hash
		From gethcommon.Address
	}
	if err := c.call(ctx, &meta, "eth_getTransactionByBlockHashAndIndex", block, gethhexutil.Uint64(index)); err != nil {
		return gethcommon.Address{}, err
	}
	if meta.Hash == (gethcommon.Hash{}) || meta.Hash != tx.Hash() {
		return gethcommon.Address{}, errors.New("wrong inclusion block/index")
	}
	return meta.From, nil
}

func toCallArg(msg *geth.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = gethhexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*gethhexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = gethhexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*gethhexutil.Big)(msg.GasPrice)
	}
	return arg
}

func toFilterArg(q geth.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = types.ToBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = types.ToBlockNumArg(q.ToBlock)
	}
	return arg, nil
}

//nolint:gocritic
func (c *Client) getBlock(ctx context.Context, method string, args ...interface{}) (*gethtypes.Block, error) {
	var raw json.RawMessage
	if err := c.call(ctx, &raw, method, args...); err != nil {
		return nil, err
	} else if len(raw) == 0 {
		return nil, geth.NotFound
	}
	// Decode header and transactions.
	var head *gethtypes.Header
	var body types.RPCBlock
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	if head.UncleHash == gethtypes.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, fmt.Errorf("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != gethtypes.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, fmt.Errorf("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == gethtypes.EmptyRootHash && len(body.Transactions) > 0 {
		return nil, fmt.Errorf("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != gethtypes.EmptyRootHash && len(body.Transactions) == 0 {
		return nil, fmt.Errorf("server returned empty transaction list but block header indicates transactions")
	}
	// Load uncles because they are not included in the block response.
	var uncles []*gethtypes.Header
	// if len(body.UncleHashes) > 0 {
	// 	uncles = make([]*gethtypes.Header, len(body.UncleHashes))
	// 	reqs := make([]rpc.BatchElem, len(body.UncleHashes))
	// 	for i := range reqs {
	// 		reqs[i] = rpc.BatchElem{
	// 			Method: "eth_getUncleByBlockHashAndIndex",
	// 			Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
	// 			Result: &uncles[i],
	// 		}
	// 	}
	// 	if err := c.BatchCallContext(ctx, reqs); err != nil {
	// 		return nil, err
	// 	}
	// 	for i := range reqs {
	// 		if reqs[i].Error != nil {
	// 			return nil, reqs[i].Error
	// 		}
	// 		if uncles[i] == nil {
	// 			return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
	// 		}
	// 	}
	// }
	// Fill the sender cache of transactions in the block.
	txs := make([]*gethtypes.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		if tx.From != nil {
			setSenderFromServer(tx.Tx, *tx.From, body.Hash)
		}
		txs[i] = tx.Tx
	}
	return gethtypes.NewBlockWithHeader(head).WithBody(txs, uncles), nil
}
