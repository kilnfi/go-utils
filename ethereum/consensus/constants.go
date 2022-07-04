package ethcl

import (
	"fmt"

	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

var (
	MainnetForkVersion = beaconcommon.Version{0x00, 0x00, 0x00, 0x00}
	PraterForkVersion  = beaconcommon.Version{0x00, 0x00, 0x10, 0x20}
	SepoliaForkVersion = beaconcommon.Version{0x90, 0x00, 0x00, 0x69}
	RopstenForkVersion = beaconcommon.Version{0x80, 0x00, 0x00, 0x69}
)

var forkVersions = map[string]beaconcommon.Version{
	"mainnet": MainnetForkVersion,
	"prater":  PraterForkVersion,
	"goerli":  PraterForkVersion, // we add goerli to facilitate correspondance with exec layer
	"sepolia": SepoliaForkVersion,
	"ropsten": RopstenForkVersion,
}

var networks = map[string]string{
	MainnetForkVersion.String(): "mainet",
	PraterForkVersion.String():  "prater",
	SepoliaForkVersion.String(): "sepolia",
	RopstenForkVersion.String(): "ropsten",
}

func ForkVersion(network string) (beaconcommon.Version, error) {
	if v, ok := forkVersions[network]; ok {
		return v, nil
	}
	return beaconcommon.Version{}, fmt.Errorf("unkown network %v", network)
}

func Network(v beaconcommon.Version) (string, error) {
	if v, ok := networks[v.String()]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unkown fork version %v", v)
}
