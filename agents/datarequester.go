package agents

import (
	"proposalsubmitters/entities"
	"proposalsubmitters/utils"
)

type DataRequester struct {
	RPCClient *utils.HttpClient
}

func (dr *DataRequester) SubmitProposal(proposal *entities.SubmitDCBProposalMeta) (*entities.DCBProposalRes, error) {
	method := utils.SubmitDCBProposalMethod
	res := &entities.DCBProposalRes{}
	err := dr.createAndSendTx(method, proposal, res)
	return res, err
}

func (dr *DataRequester) createAndSendTx(method string, meta, rpcResponse interface{}) error {
	privKey := ""
	params := []interface{}{
		privKey,
		nil,
		DefaultFee,
		-1,
		meta,
	}
	return dr.RPCClient.RPCCall(method, params, rpcResponse)
}
