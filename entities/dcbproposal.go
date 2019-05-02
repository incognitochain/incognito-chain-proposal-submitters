package entities

import (
	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
)

type DCBProposalResult struct {
	TxID    string `json:"txid"`
	ShardID byte
}

type DCBProposalResponse struct {
	RPCBaseRes
	Result DCBProposalResult `json:"result"`
}

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

type DCBBondInfo struct {
	Amount uint64
	BondID common.Hash
	Price  uint64 // average price when buying bonds

	Maturity uint64
	BuyBack  uint64 // price for selling to GOV
}
