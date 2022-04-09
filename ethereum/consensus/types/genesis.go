package types

import (
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

type Genesis struct {
	GenesisTime           beaconcommon.Timestamp `json:"genesis_time" yaml:"genesis_time"`
	GenesisValidatorsRoot beaconcommon.Root      `json:"genesis_validators_root" yaml:"genesis_validators_root"`
	GenesisForkVersion    beaconcommon.Version   `json:"genesis_fork_version" yaml:"genesis_fork_version"`
}
