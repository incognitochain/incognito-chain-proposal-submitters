package entities

import (
	"github.com/constant-money/constant-chain/blockchain"
	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
)

// DCBProposalResponse stores response of CreateAndSendSubmitDCBProposalTx rpc
type DCBProposalResponse struct {
	RPCBaseRes
	Result DCBProposalResult `json:"result"`
}

type DCBProposalResult struct {
	TxID    string `json:"txid"`
	ShardID byte
}

// DCBProposalResponse stores payload of CreateAndSendSubmitDCBProposalTx rpc
type SubmitDCBProposalMeta struct {
	DCBParams         *component.DCBParams
	ExecuteDuration   uint64
	Explanation       string
	PaymentAddress    string
	ConstitutionIndex uint32
}

// Proposal the submitted proposal data and metadata
type Proposal struct {
	Data         *SubmitDCBProposalMeta
	ProposedTxID string

	AcceptedHeight  uint64
	SubmittedHeight uint64
}

// DCBBondInfo contains amount and average price of bonds that DCB holds
type DCBBondInfo struct {
	Amount uint64
	BondID common.Hash
	Price  uint64 // average price when buying bonds

	Maturity uint64
	BuyBack  uint64 // price for selling to GOV
}

// DCBProposalInfo stores a proposal's params as well as its term and start/end block
type DCBProposalInfo struct {
	DCBParams          *component.DCBParams
	EndBlock           uint64
	ConstitutionIndex  uint32
	StartedBlockHeight uint64
}

// StabilityInfoResponse stores response of GetCurrentStabilityInfo rpc
type StabilityInfoResponse struct {
	RPCBaseRes
	Result *blockchain.StabilityInfo
}

// BeaconBestStateResponse stores response of GetBeaconBestState rpc
type BeaconBestStateResponse struct {
	RPCBaseRes
	Result *BeaconBestStateResult
}

type BeaconBestStateResult struct {
	BeaconHeight uint64
}
