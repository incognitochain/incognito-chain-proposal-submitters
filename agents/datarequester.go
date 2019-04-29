package agents

import (
	"fmt"
	"proposalsubmitters/entities"
	"proposalsubmitters/utils"
)

type DataRequester struct {
	RPCClient *utils.HttpClient

	privKey string
}

func (dr *DataRequester) SubmitProposal(proposal *entities.SubmitDCBProposalMeta) (*entities.DCBProposalResponse, error) {
	method := utils.SubmitDCBProposalMethod
	res := &entities.DCBProposalResponse{}
	err := dr.createAndSendTx(method, proposal, res)
	fmt.Printf("res: %+v\n", res)
	fmt.Printf("res rpcerr: %+v\n", res.RPCError)
	return res, err
}

func (dr *DataRequester) createAndSendTx(method string, meta, rpcResponse interface{}) error {
	params := []interface{}{
		dr.privKey,
		nil,
		DefaultFee,
		-1,
		meta,
	}
	return dr.RPCClient.RPCCall(method, params, rpcResponse)
}
