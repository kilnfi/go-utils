package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// ErrorMsg is a struct allowing to encode/decode a JSON-RPC response body
type ErrorMsg struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

func (err ErrorMsg) Error() string {
	b, _ := json.Marshal(err)
	return fmt.Sprintf("JSON-RPC: %v", string(b))
}
