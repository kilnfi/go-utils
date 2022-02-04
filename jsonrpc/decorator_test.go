package jsonrpc_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/skillz-blockchain/go-utils/jsonrpc"
	jsonrpctestutils "github.com/skillz-blockchain/go-utils/jsonrpc/testutils"
)

func TestWithVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpctestutils.NewMockClient(ctrl)
	c := jsonrpc.WithVersion("2.0")(mockCli)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpctestutils.HasVersion("2.0"),
		gomock.Any())
	c.Call(context.Background(), &jsonrpc.Request{}, nil)
}

func TestWithIncrementalID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpctestutils.NewMockClient(ctrl)
	c := jsonrpc.WithIncrementalID()(mockCli)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpctestutils.HasID(uint32(0)),
		gomock.Any())
	c.Call(context.Background(), &jsonrpc.Request{}, nil)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpctestutils.HasID(uint32(1)),
		gomock.Any())
	c.Call(context.Background(), &jsonrpc.Request{}, nil)
}
