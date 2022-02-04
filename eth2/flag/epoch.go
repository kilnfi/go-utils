package flag

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/spf13/pflag"
)

// epochValue is a type implenting pflag.Value interface for beaconcommon.Epoch type
type epochValue struct {
	epoch *beaconcommon.Epoch
}

func (v *epochValue) Set(s string) error { return v.epoch.UnmarshalJSON([]byte(s)) }
func (v *epochValue) Type() string       { return "epoch" }
func (v *epochValue) String() string     { return v.epoch.String() }

// EpochVar registers a beaconcommon.Epoch custom flag with specified name, default value, and usage string.
// The argument p points to a beaconcommon.Epoch variable in which to store the value of the flag
func EpochVar(f *pflag.FlagSet, p *beaconcommon.Epoch, name string, value beaconcommon.Epoch, usage string) {
	v := &epochValue{p}
	*v.epoch = value
	f.Var(v, name, usage)
}

// EpochVar registers a beaconcommon.Epoch custom flag with specified name and shorthand, default value, and usage string.
// The argument p points to a beaconcommon.Epoch variable in which to store the value of the flag
func EpochVarP(f *pflag.FlagSet, p *beaconcommon.Epoch, name, shorthand string, value beaconcommon.Epoch, usage string) {
	v := &epochValue{p}
	*v.epoch = value
	f.VarP(v, name, shorthand, usage)
}
