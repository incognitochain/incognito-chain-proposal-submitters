package agents

import (
	"errors"
	"fmt"
	"proposalsubmitters/entities"
	"proposalsubmitters/utils"
)

type Fiscal1 struct {
	AgentAbs
}

func getDefaultGOVParams() *entities.GOVParams {
	oracleNetwork := &entities.OracleNetwork{
		OraclePubKeys:          []string{},
		UpdateFrequency:        10,
		Quorum:                 1,
		OracleRewardMultiplier: 1, // 0.01C
		AcceptableErrorMargin:  200,
		WrongTimesAllowed:      2,
	}
	return &entities.GOVParams{
		SalaryPerTx:      1,
		BasicSalary:      0,
		FeePerKbTx:       2,
		SellingBonds:     nil,
		SellingGOVTokens: nil,
		RefundInfo:       nil,
		OracleNetwork:    oracleNetwork,
	}
}

func buildDefaultGOVProposal() *entities.GOVProposal {
	return &entities.GOVProposal{
		GOVParams:         getDefaultGOVParams(),
		ExecuteDuration:   DefaultExecuteDuration,
		Explanation:       "",
		PaymentAddress:    utils.GetENV("SUBMITTER_PAYMENT_ADDRESS", ""),
		ConstitutionIndex: 1,
	}
}

func (f1 *Fiscal1) submitDefaultGOVProposal() (interface{}, error) {
	privKey := utils.GetENV("PRIV_KEY", "")
	params := []interface{}{
		privKey,
		nil,
		DefaultFee,
		-1,
		buildDefaultGOVProposal(),
	}
	var govProposalRes entities.GOVProposalRes
	err := f1.RPCClient.RPCCall("createandsendsubmitgovproposaltx", params, &govProposalRes)
	if err != nil {
		return nil, err
	}
	if govProposalRes.RPCError != nil {
		return nil, errors.New(govProposalRes.RPCError.Message)
	}
	return govProposalRes.Result, nil
}

func (f1 *Fiscal1) Execute() {
	fmt.Println("Fiscal1 agent is executing...")

	result, err := f1.submitDefaultGOVProposal()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Result: ", result)
}
