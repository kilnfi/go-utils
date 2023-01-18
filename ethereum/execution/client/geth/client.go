package geth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

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

// While the last release of go-ethereum can't make a call to eth_getBlockByNumber
// with finalized arg, we hook into the method to make this custom call.
// Once go-ethereum will release this:
// https://github.com/ethereum/go-ethereum/blob/c1aa1db69e74c71f251fc83cf7c120b4d0222728/ethclient/gethclient/gethclient.go#L189
// then we could simply remove this condition
// finalizedBlock, err := s.core.ElClient().BlockByNumber(ctx, big.NewInt(int64(rpc.finalizedBlockNumber)))
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	finalized := big.NewInt(int64(rpc.FinalizedBlockNumber))
	if number != nil && number.Cmp(finalized) == 0 {
		var raw json.RawMessage
		if err := c.rpcclient.CallContext(ctx, &raw, "eth_getBlockByNumber", "finalized", true); err != nil {
			return nil, err
		}
		var head *types.Header
		if err := json.Unmarshal(raw, &head); err != nil {
			return nil, err
		}
		return types.NewBlockWithHeader(head), nil // this block object is incomplete but enough for current usage.
	}

	return c.Client.BlockByNumber(ctx, number)
}
