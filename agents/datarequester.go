package agents

import (
	"fmt"
	"proposalsubmitters/entities"
	"proposalsubmitters/utils"

	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/rpcserver"
)

type DataRequester struct {
	RPCClient *utils.HttpClient

	privKey string
}

func (dr *DataRequester) SubmitProposal(proposal *entities.SubmitDCBProposalMeta) (*entities.DCBProposalResponse, error) {
	method := rpcserver.CreateAndSendSubmitDCBProposalTx
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

func (dr *DataRequester) BlockHeight() (uint64, error) {
	method := rpcserver.GetBeaconBestState
	params := []interface{}{}
	resp := &entities.BeaconBestStateResponse{}
	err := dr.RPCClient.RPCCall(method, params, resp)
	if err != nil || resp.RPCError != nil {
		return 0, entities.AggErr(err, resp.RPCError)
	}

	return resp.Result.BeaconHeight, nil
}

func (dr *DataRequester) ConstantCirculating() (uint64, error) {
	return 0, nil
}

func (dr *DataRequester) AssetPrice(assetID common.Hash) (uint64, error) {
	return 0, nil
}

func (dr *DataRequester) DCBBondPortfolio() ([]*entities.DCBBondInfo, error) {
	return nil, nil
}

func (dr *DataRequester) OngoingProposalInfo() (*entities.DCBProposalInfo, error) {
	method := rpcserver.GetCurrentStabilityInfo
	params := []interface{}{}
	resp := &entities.StabilityInfoResponse{}
	err := dr.RPCClient.RPCCall(method, params, resp)
	if err != nil || resp.RPCError != nil {
		return nil, entities.AggErr(err, resp.RPCError)
	}

	c := resp.Result.DCBConstitution
	return &entities.DCBProposalInfo{
		DCBParams:          &c.DCBParams,
		EndBlock:           c.StartedBlockHeight + c.ExecuteDuration,
		ConstitutionIndex:  c.ConstitutionIndex,
		StartedBlockHeight: c.StartedBlockHeight,
	}, nil
}
