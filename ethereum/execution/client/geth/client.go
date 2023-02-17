package geth

import (
	"context"
	"fmt"
	"math/big"
	"sync"

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
