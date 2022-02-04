package types

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

type SyncCommittees struct {
	Validators           beaconcommon.CommitteeIndices   `json:"validators"`
	ValidatorsAggregates []beaconcommon.CommitteeIndices `json:"validator_aggregates"`
}
