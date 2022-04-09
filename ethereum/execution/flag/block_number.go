package flag

import (
	"math/big"

	"github.com/spf13/pflag"

	"github.com/skillz-blockchain/go-utils/ethereum/execution/types"
)

type blockNumberValue struct {
	block **big.Int
}

func (b blockNumberValue) String() string { return types.ToBlockNumArg(*b.block) }
func (b blockNumberValue) Type() string   { return "blockNumber" }
func (b *blockNumberValue) Set(s string) error {
	bn, err := types.FromBlockNumArg(s)
	if err != nil {
		return err
	}
	*b.block = bn
	return nil
}

// BlockNumberVar registers a *big.Int blocknumber custom flag with specified name, default value, and usage string.
// The argument p points to a *big.Int variable in which to store the value of the flag
func BlockNumberVar(f *pflag.FlagSet, p **big.Int, name string, value *big.Int, usage string) {
	*p = value
	f.Var(&blockNumberValue{p}, name, usage)
}

// BlockNumberVar registers a *big.Int blocknumber custom flag with specified name and shorthand, default value, and usage string.
// The argument p points to a *big.Int variable in which to store the value of the flag
func BlockNumberVarP(f *pflag.FlagSet, b **big.Int, name, shorthand string, value *big.Int, usage string) {
	*b = value
	f.VarP(&blockNumberValue{b}, name, shorthand, usage)
}
