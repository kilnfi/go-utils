package ethel

import "fmt"

var (
	MainnetChainID = "1"
	RopstenChainID = "3"
	GoerliChainID  = "5"
	SepoliaChainID = "11155111"
)

var chainIds = map[string]string{
	"mainnet": MainnetChainID,
	"prater":  GoerliChainID, // we add prater to facilitate correspondance with consensus layer
	"goerli":  GoerliChainID,
	"sepolia": SepoliaChainID,
	"ropsten": RopstenChainID,
}

var networks = map[string]string{
	MainnetChainID: "mainnet",
	GoerliChainID:  "goerli",
	SepoliaChainID: "sepolia",
	RopstenChainID: "ropsten",
}

func ChainID(network string) (string, error) {
	if v, ok := chainIds[network]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown network %v", network)
}

func Network(id string) (string, error) {
	if v, ok := networks[id]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown chain id %v", id)
}
