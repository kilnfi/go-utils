package flag

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/spf13/pflag"
)

// GweiValue is a type implenting pflag.Value interface for beaconcommon.Gwei type
type GweiValue struct {
	Gwei *beaconcommon.Gwei
}

func (v *GweiValue) Set(s string) error { return v.Gwei.UnmarshalJSON([]byte(s)) }
func (v *GweiValue) Type() string       { return "Gwei" }
func (v *GweiValue) String() string     { return v.Gwei.String() }

// GweiVar registers a beaconcommon.Gwei custom flag with specified name, default value, and usage string.
// The argument p points to a beaconcommon.Gwei variable in which to store the value of the flag
func GweiVar(f *pflag.FlagSet, p *beaconcommon.Gwei, name string, value beaconcommon.Gwei, usage string) {
	v := &GweiValue{p}
	*v.Gwei = value
	f.Var(v, name, usage)
}

// GweiVar registers a beaconcommon.Gwei custom flag with specified name and shorthand, default value, and usage string.
// The argument p points to a beaconcommon.Gwei variable in which to store the value of the flag
func GweiVarP(f *pflag.FlagSet, p *beaconcommon.Gwei, name, shorthand string, value beaconcommon.Gwei, usage string) {
	v := &GweiValue{p}
	*v.Gwei = value
	f.VarP(v, name, shorthand, usage)
}
