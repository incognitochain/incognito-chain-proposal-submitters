package entities

import (
	"github.com/constant-money/constant-chain/blockchain/component"
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
