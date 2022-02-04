package flag

import (
	"math/big"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestBlockNumber(t *testing.T) {
	t.Run("default nil and flag unset", func(t *testing.T) {
		b := big.NewInt(0)
		flags := pflag.NewFlagSet("test", pflag.PanicOnError)
		BlockNumberVar(flags, &b, "test-block", nil, "Test usage")
		_ = flags.Parse([]string{})
		assert.Nil(t, b)
	})

	t.Run("default nil and flag set", func(t *testing.T) {
		b := big.NewInt(0)
		flags := pflag.NewFlagSet("test", pflag.PanicOnError)
		BlockNumberVar(flags, &b, "test-block", nil, "Test usage")
		_ = flags.Parse([]string{"--test-block", "0x1"})
		assert.Equal(t, big.NewInt(1), b)
	})

	t.Run("default not nil and flag set", func(t *testing.T) {
		b := big.NewInt(0)
		flags := pflag.NewFlagSet("test", pflag.PanicOnError)
		BlockNumberVar(flags, &b, "test-block", nil, "Test usage")
		_ = flags.Parse([]string{"--test-block", "0x1"})
		assert.Equal(t, big.NewInt(1), b)
	})
}
