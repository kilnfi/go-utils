package types

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

type BeaconBlockHeader struct {
	Root      beaconcommon.Root                    `json:"root"`
	Canonical bool                                 `json:"canonical"`
	Header    beaconcommon.SignedBeaconBlockHeader `json:"header"`
}
