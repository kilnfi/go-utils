//go:build !integration
// +build !integration

package flag

import (
	"testing"

	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {
	root := beaconcommon.Root{}
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RootVar(f, &root, "test-root", beaconcommon.Root{}, "test")

	// No 0x prefix arg
	err := f.Parse([]string{"--test-root=0100000000000000000000007e654d251da770a068413677967f6d3ea2fea9e4"})
	require.NoError(t, err)
	assert.Equal(t, "0x0100000000000000000000007e654d251da770a068413677967f6d3ea2fea9e4", root.String())

	// 0x prefix arg
	err = f.Parse([]string{"--test-root=0x0100000000000000000000008f654d251da770a068413677967f6d3ea2fea9e4"})
	require.NoError(t, err)
	assert.Equal(t, "0x0100000000000000000000008f654d251da770a068413677967f6d3ea2fea9e4", root.String())
}
