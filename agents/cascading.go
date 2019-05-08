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

	Proposal *entities.Proposal
	Data     *DataRequester

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
			Frequency: 12,
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

func (ca *CascadingAgent) submitDCBProposal(proposal *entities.SubmitDCBProposalMeta) (*entities.Proposal, error) {
	blockHeight, err := ca.Data.BlockHeight()
	if err != nil {
		return nil, err
	}

	resp, err := ca.Data.SubmitProposal(proposal)
	if err != nil || resp.RPCError != nil {
		return nil, entities.AggErr(err, resp.RPCError)
	}

	return &entities.Proposal{
		Data:            proposal,
		ProposedTxID:    resp.Result.TxID,
		SubmittedHeight: blockHeight,
	}, nil
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
		ExecuteDuration:   100,
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

	burnAmount := uint64(float64(Peg-price) * float64(circulation) / float64(Peg))
	if burnAmount < MinAmountToBurn {
		return nil, nil
	}

	fmt.Println("amount to burn:", burnAmount)
	if ca.NumSaleAccepted < NumSaleToTry {
		// Sell bonds to open market
		sales, errSale = buildCrowdsalesSellBond(burnAmount, price, blockHeight, bonds)
	}

	if sales == nil && ca.NumTradeAccepted < NumTradeToTry {
		// Sell bonds to GOV
		trades, errTrade = buildTradeSellBond(burnAmount, blockHeight, bonds)
	}

	if sales == nil && trades == nil {
		// Spend reserve
		spends, errSpend = buildSpendReserve(burnAmount, price, blockHeight, ca.Data)
	}

	if sales == nil && trades == nil && spends == nil {
		return nil, entities.CheckError(errSale, errTrade, errSpend)
	}
	proposal := ca.buildProposal(sales, trades, spends, nil)
	return proposal, nil
}

func (ca *CascadingAgent) buildExpandingProposal(price uint64) (*entities.SubmitDCBProposalMeta, error) {
	blockHeight, err := ca.Data.BlockHeight()
	if err != nil {
		return nil, err
	}

	var sales []component.SaleData
	var trades []*component.TradeBondWithGOV
	var errSale, errTrade error

	circulation, err := ca.Data.ConstantCirculating()
	if err != nil {
		return nil, err
	}

	mintAmount := uint64(float64(price-Peg) * float64(circulation) / float64(Peg))
	if mintAmount < MinAmountToMint {
		return nil, nil
	}

	fmt.Println("amount to mint:", mintAmount)
	if ca.NumSaleAccepted < NumSaleToTry {
		// Buy bonds from open market
		sales, errSale = buildCrowdsalesBuyBond(mintAmount, price, blockHeight, ca.Data)
	}

	if sales == nil && ca.NumTradeAccepted < NumTradeToTry {
		// Buy bonds from GOV
		trades, errTrade = buildTradeBuyBond(mintAmount, blockHeight, ca.Data)
	}

	if sales == nil && trades == nil {
		return nil, entities.CheckError(errSale, errTrade)
	}
	proposal := ca.buildProposal(sales, trades, nil, nil)
	return proposal, nil
}

func (ca *CascadingAgent) evaluatingProposal() (bool, error) {
	blockHeight, err := ca.Data.BlockHeight()
	if err != nil {
		return false, err
	}

	// Check current proposal and if we need to submit a new one
	info, err := ca.Data.OngoingProposalInfo()
	if err != nil {
		return false, err
	}

	evaluationEndBlock := (info.StartedBlockHeight + info.EndBlock) / 2
	if evaluationEndBlock > blockHeight {
		return true, nil // Current proposal was submitted recently
	}

	if ca.Proposal == nil {
		return false, nil // No proposal submitted
	}

	if ca.Proposal.Data.ConstitutionIndex == info.ConstitutionIndex {
		return true, nil // Voting submitted proposal
	}

	// Increase count if new proposal has sale, trade or reserve spending
	// TODO(@0xbunyip): check content of sales/trades, whether it matches submitting proposal
	if len(info.DCBParams.ListSaleData) > 0 {
		ca.NumSaleAccepted += 1
	}
	if len(info.DCBParams.TradeBonds) > 0 {
		ca.NumTradeAccepted += 1
	}

	ca.Proposal = nil // New term
	return false, nil
}

func (ca *CascadingAgent) buildHoldingPegProposal(price uint64) (*entities.SubmitDCBProposalMeta, error) {
	ca.NumSaleAccepted = 0
	ca.NumTradeAccepted = 0
	return nil, nil
}

func (ca *CascadingAgent) Execute() {
	fmt.Println("CascadingAgent agent is executing...")
	if evaluating, err := ca.evaluatingProposal(); err != nil {
		fmt.Println(err)
		return
	} else if evaluating {
		fmt.Println("Evaluating ongoing proposal...")
		return
	}

	price, err := ca.Data.AssetPrice(common.ConstantID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var proposal *entities.SubmitDCBProposalMeta
	if price < Peg-PriceFloor {
		// Price is below peg, reduce supply
		proposal, err = ca.buildContractingProposal(price)
	} else if price > Peg+PriceFloor {
		// Price is above peg, increase supply
		proposal, err = ca.buildExpandingProposal(price)
	} else {
		proposal, err = ca.buildHoldingPegProposal(price)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	// Submit proposal and save TxID
	if proposal == nil {
		fmt.Println("Assets unavailable, no proposal created")
		return
	}

	p, err := ca.submitDCBProposal(proposal)
	if err != nil {
		fmt.Println(err)
		return
	}
	ca.Proposal = p
}
