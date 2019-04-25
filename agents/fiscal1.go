package agents

import (
	"fmt"
	"proposalsubmitters/entities"
)

type Fiscal1 struct {
	AgentAbs
}

func (f1 *Fiscal1) Execute() {
	fmt.Println("Fiscal1 agent is executing...")
	var curSellingGOVRes entities.RPCCurrentSellingGOVTokensRes
	err := f1.RPCClient.RPCCall("getcurrentsellinggovtokens", nil, &curSellingGOVRes)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		return
	}
	if curSellingGOVRes.RPCError != nil {
		fmt.Errorf("Error: %s", curSellingGOVRes.RPCError.Message)
		return
	}
	fmt.Println("GOVTokenID: ", curSellingGOVRes.Result.GOVTokenID)
	fmt.Println("StartSellingAt: ", curSellingGOVRes.Result.StartSellingAt)
}
