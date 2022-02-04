package testutils

import (
	"fmt"

	"github.com/golang/mock/gomock"

	jsonrpc "github.com/skillz-blockchain/go-utils/jsonrpc"
)

type matcher struct {
	match func(req *jsonrpc.Request) bool
	msg   string
}

func (m *matcher) Matches(x interface{}) bool {
	req, ok := x.(*jsonrpc.Request)
	if !ok {
		return false
	}

	return m.match(req)
}

func (m *matcher) String() string {
	return m.msg
}

func HasVersion(v string) gomock.Matcher {
	return &matcher{
		match: func(req *jsonrpc.Request) bool { return req.Version == v },
		msg:   fmt.Sprintf("Request should have version %q", v),
	}
}

func HasID(id interface{}) gomock.Matcher {
	return &matcher{
		match: func(req *jsonrpc.Request) bool { return req.ID == id },
		msg:   fmt.Sprintf("Request should have version %v", id),
	}
}
