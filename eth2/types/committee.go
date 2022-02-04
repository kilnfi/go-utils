package types

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

type Committee struct {
	Slot       beaconcommon.Slot             `json:"slot"`
	Index      beaconcommon.CommitteeIndex   `json:"index"`
	Validators beaconcommon.CommitteeIndices `json:"validators"`
}
