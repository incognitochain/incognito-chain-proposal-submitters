package agents

import (
	"fmt"
	"proposalsubmitters/entities"
	"proposalsubmitters/utils"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
)

type CascadingAgent struct {
	AgentAbs

	ProposedTxID string
	Data         *DataRequester

	NumSaleAccepted  int
	NumTradeAccepted int

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

func (ca *CascadingAgent) submitDCBProposal(proposal *entities.SubmitDCBProposalMeta) (*entities.DCBProposalResponse, error) {
	return ca.Data.SubmitProposal(proposal)
}

func (ca *CascadingAgent) buildProposal(
	sales []component.SaleData,
	trades []*component.TradeBondWithGOV,
	spendReserveData map[common.Hash]*component.SpendReserveData,
	raiseReserveData map[common.Hash]*component.RaiseReserveData,
) *entities.SubmitDCBProposalMeta {
	if sales == nil {
		sales = []component.SaleData{}
	}
	if trades == nil {
		trades = []*component.TradeBondWithGOV{}
	}
	return &entities.SubmitDCBProposalMeta{
		DCBParams: &component.DCBParams{
			ListSaleData:             sales,
			TradeBonds:               trades,
			MinLoanResponseRequire:   0,
			MinCMBApprovalRequire:    0,
			LateWithdrawResponseFine: 0,
			RaiseReserveData:         raiseReserveData,
			SpendReserveData:         spendReserveData,
			DividendAmount:           0,
			ListLoanParams:           []component.LoanParams{},
		},
		ExecuteDuration:   1000,
		Explanation:       "Bot proposal",
		PaymentAddress:    ca.payment,
		ConstitutionIndex: 1, // TODO(@0xbunyip) update CI
	}
}

func (ca *CascadingAgent) buildContractingProposal(price uint64) (*entities.SubmitDCBProposalMeta, error) {
	bonds, err := ca.Data.DCBBondPortfolio()
	if err != nil {
		return nil, err
	}

	blockHeight, err := ca.Data.BlockHeight()
	if err != nil {
		return nil, err
	}

	var sales []component.SaleData
	var trades []*component.TradeBondWithGOV
	var spends map[common.Hash]*component.SpendReserveData
	var errSale, errTrade, errSpend error

	circulation, err := ca.Data.ConstantCirculating()
	if err != nil {
		return nil, err
	}

	burnAmount := uint64(float64(price-Peg) * float64(circulation))
	if ca.NumSaleAccepted == 0 {
		// Sell bonds to open market
		sales, errSale = buildCrowdsalesSellBond(burnAmount, price, blockHeight, bonds)
	}

	if sales == nil && ca.NumTradeAccepted == 0 {
		// Sell bonds to GOV
		trades, errTrade = buildTradeSellBond(burnAmount, blockHeight, bonds)
	}

	if sales == nil && trades == nil {
		// Spend reserve
		spends, errSpend = buildSpendReserve(burnAmount, price, blockHeight, ca.Data)
	}

	if sales == nil && trades == nil && spends == nil {
		return nil, common.CheckError(errSale, errTrade, errSpend)
	}
	proposal := ca.buildProposal(sales, trades, spends, nil)
	return proposal, nil
}

func (ca *CascadingAgent) buildExpandingProposal(price uint64) (*entities.SubmitDCBProposalMeta, error) {
	return nil, nil
}

func (ca *CascadingAgent) Execute() {
	fmt.Println("CascadingAgent agent is executing...")
	// If a proposal was submitted, wait for voting
	if len(ca.ProposedTxID) > 0 {
	}

	// Wait if our proposal has just been accepted

	price, err := ca.Data.AssetPrice(common.ConstantID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var proposal *entities.SubmitDCBProposalMeta
	if price > Peg+PriceCeiling {
		// Price is above peg, reduce supply
		proposal, err = ca.buildContractingProposal(price)

	} else if price < Peg+PriceFloor {
		// Price is below peg, increase supply
		proposal, err = ca.buildExpandingProposal(price)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	// Submit proposal and save TxID
	res, err := ca.submitDCBProposal(proposal)
	if err != nil || res.RPCError != nil {
		fmt.Println(err, res.RPCError)
		return
	}
	ca.ProposedTxID = res.Result.TxID
}
