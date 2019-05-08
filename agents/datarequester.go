package agents

import (
	"fmt"
	"proposalsubmitters/entities"
	"proposalsubmitters/utils"

	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/rpcserver"
	"github.com/constant-money/constant-chain/rpcserver/jsonresult"
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
	method := rpcserver.GetConstantCirculating
	params := []interface{}{}
	resp := &entities.ConstantCirculatingResponse{}
	err := dr.RPCClient.RPCCall(method, params, resp)
	if err != nil || resp.RPCError != nil {
		return 0, entities.AggErr(err, resp.RPCError)
	}
	return resp.Result.Total, nil
}

func (dr *DataRequester) AssetPrice(assetID common.Hash) (uint64, error) {
	method := rpcserver.GetAssetPrice
	params := []interface{}{assetID.String()}
	resp := &entities.AssetPriceResponse{}
	err := dr.RPCClient.RPCCall(method, params, resp)
	if err != nil || resp.RPCError != nil {
		return 0, entities.AggErr(err, resp.RPCError)
	}

	return resp.Result, nil
}

func (dr *DataRequester) BondsCirculating() ([]*entities.DCBBondInfo, error) {
	return nil, nil
}

func (dr *DataRequester) CurrentSellingBond() (*entities.DCBBondInfo, error) {
	return nil, nil
}

func (dr *DataRequester) bondTypes() ([]*jsonresult.GetBondTypeResultItem, error) {
	method := rpcserver.GetBondTypes
	params := []interface{}{}
	resp := &entities.BondTypesResponse{}
	err := dr.RPCClient.RPCCall(method, params, resp)
	if err != nil || resp.RPCError != nil {
		return nil, entities.AggErr(err, resp.RPCError)
	}

	result := []*jsonresult.GetBondTypeResultItem{}
	for _, b := range resp.Result.BondTypes {
		result = append(result, &b)
	}
	return result, nil
}

func (dr *DataRequester) DCBBondPortfolio() ([]*entities.DCBBondInfo, error) {
	portfolio, err := dr.buildDCBBondPortfolio()
	if err != nil || len(portfolio) == 0 {
		return nil, err
	}

	bonds, err := dr.bondTypes()
	if err != nil {
		return nil, err
	}

	// Set price and maturity
	for _, p := range portfolio {
		for _, b := range bonds {
			bondID, err := common.NewHashFromStr(b.BondID)
			if err != nil {
				continue
			}

			if p.BondID.IsEqual(bondID) {
				p.Price = b.BuyPrice // TODO(@0xbunyip): set average buying price
				p.Maturity = b.Maturity + b.StartSellingAt
				p.BuyBack = b.BuyBackPrice
			}
		}
	}
	return portfolio, nil
}

func (dr *DataRequester) buildDCBBondPortfolio() ([]*entities.DCBBondInfo, error) {
	method := rpcserver.GetListCustomTokenBalance
	dcbPayment := "1Uv46Pu4pqBvxCcPw7MXhHfiAD5Rmi2xgEE7XB6eQurFAt4vSYvfyGn3uMMB1xnXDq9nRTPeiAZv5gRFCBDroRNsXJF1sxPSjNQtivuHk"
	params := []interface{}{dcbPayment}
	resp := &entities.ListCustomTokenBalanceResponse{}
	err := dr.RPCClient.RPCCall(method, params, resp)
	if err != nil || resp.RPCError != nil {
		return nil, entities.AggErr(err, resp.RPCError)
	}

	portfolio := []*entities.DCBBondInfo{}
	for _, t := range resp.Result.ListCustomTokenBalance {
		tid, err := common.NewHashFromStr(t.TokenID)
		if err != nil || !common.IsBondAsset(tid) {
			continue
		}

		portfolio = append(portfolio, &entities.DCBBondInfo{
			Amount: t.Amount,
			BondID: *tid,
		})
	}
	return portfolio, nil
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
