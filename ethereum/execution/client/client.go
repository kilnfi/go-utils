package client

import (
	"context"
	"fmt"
	"math/big"

	geth "github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"

	"github.com/kilnfi/go-utils/common/interfaces"
	"github.com/kilnfi/go-utils/ethereum/execution/types"
	"github.com/kilnfi/go-utils/net/jsonrpc"
	jsonrpchttp "github.com/kilnfi/go-utils/net/jsonrpc/http"
)

// Client provides methods to interface with a JSON-RPC Ethereum 1.0 node
type Client struct {
	client jsonrpc.Client
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

// ChainID returns chain id
func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	res := new(gethhexutil.Big)
	err := c.call(ctx, res, "eth_chainId")
	if err != nil {
		return nil, err
	}

	return (*big.Int)(res), nil
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

// HeaderByNumber returns header a given block number
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
//nolint:gocritic
func (c *Client) CallContract(ctx context.Context, msg geth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	res := new(gethhexutil.Bytes)
	err := c.call(ctx, res, "eth_call", toCallArg(&msg), types.ToBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	return []byte(*res), nil
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
