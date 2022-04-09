package flag

import (
	"math/big"

	"github.com/skillz-blockchain/go-utils/ethereum/execution/types"

	"github.com/spf13/pflag"
)

type bigIntValue struct {
	block **big.Int
}

func (b bigIntValue) String() string { return types.EncodeBig(*b.block) }
func (b bigIntValue) Type() string   { return "bigInt" }
func (b *bigIntValue) Set(s string) error {
	bn, err := types.DecodeBig(s)
	if err != nil {
		return err
	}
	*b.block = bn
	return nil
}

// BigIntVar register a *big.Int custom flag with specified name, default value, and usage string.
// The argument p points to a *big.Int variable in which to store the value of the flag
func BigIntVar(f *pflag.FlagSet, p **big.Int, name string, value *big.Int, usage string) {
	*p = value
	f.Var(&bigIntValue{p}, name, usage)
}

// BigIntVar registers a *big.Int custom flag with specified name and shorthand, default value, and usage string.
// The argument p points to a *big.Int variable in which to store the value of the flag
func BigIntVarP(f *pflag.FlagSet, p **big.Int, name, shorthand string, value *big.Int, usage string) {
	*p = value
	f.VarP(&bigIntValue{p}, name, shorthand, usage)
}
