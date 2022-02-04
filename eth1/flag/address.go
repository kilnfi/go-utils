package flag

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"
)

// addressValue is a type implenting pflag.Value interface for gethcommon.Address type
type addressValue struct {
	addr *gethcommon.Address
}

func (v *addressValue) Set(s string) error {
	if !gethcommon.IsHexAddress(s) {
		return fmt.Errorf("invalid Ethereum address %q", s)
	}

	v.addr.SetBytes(gethcommon.HexToAddress(s).Bytes())

	return nil
}

func (v *addressValue) Type() string   { return "address" }
func (v *addressValue) String() string { return v.addr.String() }

// AddressVar registers a gethcommon.Address custom flag with specified name, default value, and usage string.
// The argument p points to a gethcommon.Address variable in which to store the value of the flag
func AddressVar(f *pflag.FlagSet, p *gethcommon.Address, name string, value gethcommon.Address, usage string) {
	p.SetBytes(value.Bytes())
	f.Var(&addressValue{p}, name, usage)
}

// AddressVar registers a gethcommon.Address custom flag with specified name, and shorthand, default value, and usage string.
// The argument p points to a gethcommon.Address variable in which to store the value of the flag
func AddressVarP(f *pflag.FlagSet, p *gethcommon.Address, name, shorthand string, value gethcommon.Address, usage string) {
	p.SetBytes(value.Bytes())
	f.VarP(&addressValue{p}, name, shorthand, usage)
}
