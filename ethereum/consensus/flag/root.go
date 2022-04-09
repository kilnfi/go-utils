package flag

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/spf13/pflag"
)

// RootValue is a type implenting pflag.Value interface for beaconcommon.Root type
type RootValue struct {
	Root *beaconcommon.Root
}

func (v *RootValue) Set(s string) error { return v.Root.UnmarshalText([]byte(s)) }
func (v *RootValue) Type() string       { return "Root" }
func (v *RootValue) String() string     { return v.Root.String() }

// RootVar registers a beaconcommon.Root custom flag with specified name, default value, and usage string.
// The argument p points to a beaconcommon.Root variable in which to store the value of the flag
func RootVar(f *pflag.FlagSet, p *beaconcommon.Root, name string, value beaconcommon.Root, usage string) {
	v := &RootValue{p}
	*v.Root = value
	f.Var(v, name, usage)
}

// RootVar registers a beaconcommon.Root custom flag with specified name and shorthand, default value, and usage string.
// The argument p points to a beaconcommon.Root variable in which to store the value of the flag
func RootVarP(f *pflag.FlagSet, p *beaconcommon.Root, name, shorthand string, value beaconcommon.Root, usage string) {
	v := &RootValue{p}
	*v.Root = value
	f.VarP(v, name, shorthand, usage)
}
