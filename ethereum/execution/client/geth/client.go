package geth

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kilnfi/go-utils/ethereum/execution/client"
)

// Ensure Client interface is fully implemented
var _ client.Client = (*Client)(nil)

// Wrapper for the go-ethereum client
type Client struct {
	*ethclient.Client

	address   string
	rpcclient *rpc.Client

	chainID *big.Int
	mu      sync.Mutex
}

func NewClient(address string) *Client {
	return &Client{
		address: address,
	}
}

func (c *Client) Init(ctx context.Context) error {
	rpcClient, err := rpc.Dial(c.address)
	if err != nil {
		return fmt.Errorf("failed to connect execution layer: %w", err)
	}
	c.rpcclient = rpcClient
	c.Client = ethclient.NewClient(rpcClient)
	return nil
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.chainID != nil {
		return c.chainID, nil
	}
	id, err := c.Client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	c.chainID = id
	return id, nil
}

func (c *Client) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	logs, err := c.Client.FilterLogs(ctx, query)
	if err == nil {
		return logs, nil
	}

	type jsonError interface {
		Error() string
		ErrorCode() int
		ErrorData() interface{}
	}

	jsonErr, ok := err.(jsonError)
	if !ok {
		return logs, err
	}

	// In some case we want to retry the request with a smaller block range
	switch jsonErr.ErrorCode() {
	case
		-32005, // LimitExceededError
		-32002: // RPC timeout
		return c.splitFilterLogs(ctx, query)
	}

	return logs, err
}

func (c *Client) splitFilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	var (
		diff        = new(big.Int).Sub(query.ToBlock, query.FromBlock)
		offset      = new(big.Int).Div(diff, big.NewInt(2))
		middleBlock = new(big.Int).Add(query.FromBlock, offset)
		threshold   = big.NewInt(1000)
	)

	if diff.Cmp(threshold) < 0 {
		return nil, fmt.Errorf("failed to FilterLogs: %w", context.DeadlineExceeded)
	}

	query1 := ethereum.FilterQuery{
		BlockHash: query.BlockHash,
		FromBlock: query.FromBlock,
		ToBlock:   middleBlock,
		Addresses: query.Addresses,
		Topics:    query.Topics,
	}
	logs1, err := c.FilterLogs(ctx, query1)
	if err != nil {
		return logs1, fmt.Errorf("failed to FilterLogs: %w", err)
	}

	query2 := ethereum.FilterQuery{
		BlockHash: query.BlockHash,
		FromBlock: new(big.Int).Add(middleBlock, big.NewInt(1)),
		ToBlock:   query.ToBlock,
		Addresses: query.Addresses,
		Topics:    query.Topics,
	}
	logs2, err := c.FilterLogs(ctx, query2)
	if err != nil {
		return logs2, fmt.Errorf("failed to FilterLogs: %w", err)
	}

	return append(logs1, logs2...), nil
}
