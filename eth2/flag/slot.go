package flag

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/spf13/pflag"
)

// slotValue is a type implenting pflag.Value interface for beaconcommon.Epoch type
type slotValue struct {
	slot *beaconcommon.Slot
}

func (v *slotValue) Set(s string) error { return v.slot.UnmarshalJSON([]byte(s)) }
func (v *slotValue) Type() string       { return "slot" }
func (v *slotValue) String() string     { return v.slot.String() }

// SlotVar registers a beaconcommon.Slot custom flag with specified name, default value, and usage string.
// The argument p points to a beaconcommon.Slot variable in which to store the value of the flag
func SlotVar(f *pflag.FlagSet, p *beaconcommon.Slot, name string, value beaconcommon.Slot, usage string) {
	v := &slotValue{p}
	*v.slot = value
	f.Var(v, name, usage)
}

// SlotVar registers a beaconcommon.Slot custom flag with specified name and shorthand, default value, and usage string.
// The argument p points to a beaconcommon.Slot variable in which to store the value of the flag
func SlotVarP(f *pflag.FlagSet, p *beaconcommon.Slot, name, shorthand string, value beaconcommon.Slot, usage string) {
	v := &slotValue{p}
	*v.slot = value
	f.VarP(v, name, shorthand, usage)
}
