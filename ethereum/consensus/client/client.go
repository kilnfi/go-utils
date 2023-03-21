package client

import (
	"context"

	"github.com/kilnfi/go-utils/ethereum/consensus/types"

	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	beaconphase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// Client interface to an Ethereum 2.0 node

// Note:

// For every method receiving stateID argument, stateID can be one of:
// - "head" (canonical head in node's view)
// - "genesis"
// - "finalized"
// - "justified"
// - <slot>
// - <hex encoded state root (with 0x prefix)>.

// For every method receiving validatorID argument, validatorID can be one of:
// - <hex encoded public key (with 0x prefix)>
// - <validator index>

// For every method receiving blockID argument, blockID can be one of:
// - "head" (canonical head in node's view)
// - "genesis"
// - "finalized"
// - <slot>
// - <hex encoded block root (with 0x prefix)>.
//
//go:generate mockgen -source client.go -destination mock/client.go -package mock client
type Client interface {
	BeaconClient
	NodeClient
	ConfigClient
}

type BeaconClient interface {
	// GetGenesis returns details of the chain's genesis
	GetGenesis(ctx context.Context) (*types.Genesis, error)

	// GetStateRoot calculates HashTreeRoot of the state for the given stateID
	GetStateRoot(ctx context.Context, stateID string) (*beaconcommon.Root, error)

	// GetStateFork returns Fork object for state with given stateID
	GetStateFork(ctx context.Context, stateID string) (*beaconcommon.Fork, error)

	// GetStateFinalityCheckpoints returns finality checkpoints for state with given stateID
	// In case finality is not yet achieved returns epoch 0 and ZERO_HASH as root.
	GetStateFinalityCheckpoints(ctx context.Context, stateID string) (*types.StateFinalityCheckpoints, error)

	// GetValidators returns list of validators
	// Set validatorsIDs and/or statuses to filter result (if empty no filter is applied)
	GetValidators(ctx context.Context, stateID string, validatorIDs, statuses []string) ([]*types.Validator, error)

	// GetValidator returns validator specified by stateID and validatorID
	GetValidator(ctx context.Context, stateID, validatorID string) (*types.Validator, error)

	// GetValidatorBalances returns list of validator balances.
	// Set validatorsIDs to filter validator result (if empty no filter is applied)
	GetValidatorBalances(ctx context.Context, stateID string, validatorIDs []string) ([]*types.ValidatorBalance, error)

	// GetCommittees returns the committees for the given state.
	// Set epoch and/or index and/or slot to filter result (if nil no filter is applied)
	GetCommittees(ctx context.Context, stateID string, epoch *beaconcommon.Epoch, index *beaconcommon.CommitteeIndex, slot *beaconcommon.Slot) ([]*types.Committee, error)

	// GetSyncCommittees returns the sync committees for given stateID
	// Set epoch to filter result (if nil no filter is applied)
	GetSyncCommittees(ctx context.Context, stateID string, epoch *beaconcommon.Epoch) (*types.SyncCommittees, error)

	// GetBlockHeaders return block headers
	// Set slot and/or parentRoot to filter result (if nil no filter is applied)
	GetBlockHeaders(ctx context.Context, slot *beaconcommon.Slot, parentRoot *beaconcommon.Root) ([]*types.BeaconBlockHeader, error)

	// GetBlockHeader returns block header for given blockID
	GetBlockHeader(ctx context.Context, blockID string) (*types.BeaconBlockHeader, error)

	// GetBlock returns block details for given block id.
	GetBlock(ctx context.Context, blockID string) (*bellatrix.SignedBeaconBlock, error)

	// GetBlockRoot returns hashTreeRoot of block
	GetBlockRoot(ctx context.Context, blockID string) (*beaconcommon.Root, error)

	// GetBlockAttestations returns attestations included in requested block with given blockID
	GetBlockAttestations(ctx context.Context, blockID string) (beaconphase0.Attestations, error)

	// GetAttestations returns attestations known by the node but not necessarily incorporated into any block.
	GetAttestations(ctx context.Context) (beaconphase0.Attestations, error)

	// GetAttesterSlashings returns attester slashings known by the node but not necessarily incorporated into any block.
	GetAttesterSlashings(ctx context.Context) (beaconphase0.AttesterSlashings, error)

	// GetProposerSlashings returns proposer slashings known by the node but not necessarily incorporated into any block.
	GetProposerSlashings(ctx context.Context) (beaconphase0.ProposerSlashings, error)

	// GetVoluntaryExits returns voluntary exits known by the node but not necessarily incorporated into any block.
	GetVoluntaryExits(ctx context.Context) (beaconphase0.VoluntaryExits, error)
}

type NodeClient interface {
	// GetNodeVersion returns node's version contains informations about the node processing the request

	// Example: teku/v0.12.6-dev-994997f8/osx-x86_64/adoptopenjdk-java-11
	GetNodeVersion(ctx context.Context) (string, error)
}

type ConfigClient interface {
	// GetSpec returns Ethreum 2.0 specifications configuration used on the node.
	GetSpec(ctx context.Context) (*beaconcommon.Spec, error)
}
