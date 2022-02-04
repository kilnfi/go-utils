package types

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

type StateFinalityCheckpoints struct {
	PreviousJustifiedCheckpoint beaconcommon.Checkpoint `json:"previous_justified"`
	CurrentJustifiedCheckpoint  beaconcommon.Checkpoint `json:"current_justified"`
	FinalizedCheckpoint         beaconcommon.Checkpoint `json:"finalized"`
}
