package jsonrpc

import (
	"context"
	"sync/atomic"
)

type ClientDecorator func(Client) Client

// WithVersion automatically set JSON-RPC request version
func WithVersion(v string) ClientDecorator {
	return func(c Client) Client {
		return ClientFunc(func(ctx context.Context, req *Request, res interface{}) error {
			req.Version = v
			return c.Call(ctx, req, res)
		})
	}
}

// WithIncrementalID automatically increments JSON-RPC request ID
func WithIncrementalID() ClientDecorator {
	var idCounter uint32
	return func(c Client) Client {
		return ClientFunc(func(ctx context.Context, req *Request, res interface{}) error {
			req.ID = atomic.AddUint32(&idCounter, 1) - 1
			return c.Call(ctx, req, res)
		})
	}
}
