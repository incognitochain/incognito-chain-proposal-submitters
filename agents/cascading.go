package agents

import (
	"encoding/json"
	"fmt"
	"proposalsubmitters/entities"
	"proposalsubmitters/utils"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
)

type CascadingAgent struct {
	AgentAbs

	ProposedTxID *common.Hash
	Data         *DataRequester

	privKey string
	payment string
}

func NewCascadingAgent(rpcClient *utils.HttpClient) *CascadingAgent {
	return &CascadingAgent{
		AgentAbs: AgentAbs{
			ID:        2,
			Name:      "cascading agent 1",
			Frequency: 60,
			Quit:      make(chan bool),
			RPCClient: rpcClient,
		},
		Data: &DataRequester{
			RPCClient: rpcClient,
			privKey:   utils.GetENV("DCB_AGENT_PRIVKEY", ""),
		},
		privKey: utils.GetENV("DCB_AGENT_PRIVKEY", ""),
		payment: utils.GetENV("DCB_AGENT_PAYMENT", ""),
	}
}

func (ca *CascadingAgent) defaultSubmitDCBProposalMeta() *entities.SubmitDCBProposalMeta {
	return &entities.SubmitDCBProposalMeta{
		DCBParams: &component.DCBParams{
			ListSaleData:             []component.SaleData{},
			TradeBonds:               []*component.TradeBondWithGOV{},
			MinLoanResponseRequire:   0,
			MinCMBApprovalRequire:    0,
			LateWithdrawResponseFine: 0,
			RaiseReserveData:         nil,
			SpendReserveData:         nil,
			DividendAmount:           0,
			ListLoanParams:           []component.LoanParams{},
		},
		ExecuteDuration:   100,
		Explanation:       "Default DCB proposal",
		PaymentAddress:    ca.payment,
		ConstitutionIndex: 1,
	}
}

func (ca *CascadingAgent) SubmitDCBProposal(proposal *entities.SubmitDCBProposalMeta) (*entities.DCBProposalResponse, error) {
	return ca.Data.SubmitProposal(proposal)
}

func (ca *CascadingAgent) Execute() {
	fmt.Println("CascadingAgent agent is executing...")
	proposal := ca.defaultSubmitDCBProposalMeta()
	a, err := json.Marshal(proposal)
	fmt.Printf("%v %+v\n", err, string(a))

	res, err := ca.SubmitDCBProposal(proposal)
	fmt.Println(res, err)
}
