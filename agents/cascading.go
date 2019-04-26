package agents

import (
	"encoding/json"
	"fmt"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/privacy"
)

type CascadingAgent struct {
	AgentAbs
}

type SubmitDCBProposalMeta struct {
	DCBParams *component.DCBParams
	*component.SubmitProposalInfo
}

func (ca *CascadingAgent) defaultSubmitDCBProposalMeta() *SubmitDCBProposalMeta {
	return &SubmitDCBProposalMeta{
		DCBParams: &component.DCBParams{
			ListSaleData:             nil,
			TradeBonds:               nil,
			MinLoanResponseRequire:   0,
			MinCMBApprovalRequire:    0,
			LateWithdrawResponseFine: 0,
			RaiseReserveData:         nil,
			SpendReserveData:         nil,
			DividendAmount:           0,
			ListLoanParams:           nil,
		},
		SubmitProposalInfo: &component.SubmitProposalInfo{
			ExecuteDuration:   100,
			Explanation:       "Default DCB proposal",
			PaymentAddress:    privacy.PaymentAddress{},
			ConstitutionIndex: 1,
		},
	}
}

func (ca *CascadingAgent) Execute() {
	fmt.Println("CascadingAgent agent is executing...")
	proposal := ca.defaultSubmitDCBProposalMeta()
	a, err := json.Marshal(proposal)
	fmt.Printf("%v %+v\n", err, string(a))
}
