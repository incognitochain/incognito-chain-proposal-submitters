package entities

import (
	"github.com/constant-money/constant-chain/blockchain/component"
)

type DCBProposalRes struct {
	RPCBaseRes
	Result interface{}
}

type SubmitDCBProposalMeta struct {
	DCBParams         *component.DCBParams
	ExecuteDuration   uint64
	Explanation       string
	PaymentAddress    string
	ConstitutionIndex uint32
}
