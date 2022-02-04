package jsonrpc

import (
	"encoding/json"
)

var null = json.RawMessage("null")

// Request allows to manipulate a JSON-RPC request
type Request struct {
	Version string
	Method  string
	ID      interface{}
	Params  interface{}
}

// requestMsg is a struct allowing to encode/decode a JSON-RPC request body
type requestMsg struct {
	Version string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params,omitempty"`
	ID      *json.RawMessage `json:"id,omitempty"`
}

// MarshalJSON
func (msg *Request) MarshalJSON() ([]byte, error) {
	raw := new(requestMsg)

	raw.Version = msg.Version
	raw.Method = msg.Method

	raw.ID = new(json.RawMessage)
	if msg.ID != nil {
		b, err := json.Marshal(msg.ID)
		if err != nil {
			return nil, err
		}

		*raw.ID = b
	} else {
		copy(*raw.ID, null)
	}

	raw.Params = new(json.RawMessage)
	if msg.Params != nil {
		b, err := json.Marshal(msg.Params)
		if err != nil {
			return nil, err
		}
		*raw.Params = b
	} else {
		copy(*raw.Params, null)
	}

	return json.Marshal(raw)
}
