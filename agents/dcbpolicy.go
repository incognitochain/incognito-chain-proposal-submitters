package agents

import (
	"proposalsubmitters/entities"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
)

func buildCrowdsalesSellBond(
	burnAmount uint64,
	constantPrice uint64,
	blockHeight uint64,
	bonds []*entities.DCBBondInfo,
) ([]component.SaleData, error) {
	// Check if there's a bond that can reduce burnAmount of Constant
	var bondToSell *entities.DCBBondInfo
	for _, b := range bonds {
		// TODO(0xbunyip): pick price between b.Price and b.BuyBackPrice
		if b.Price > 0 && b.Price*b.Amount >= burnAmount {
			bondToSell = b
			break
		}
	}
	if bondToSell == nil {
		return nil, nil
	}

	sellingAmount := 1 + (burnAmount-1)/bondToSell.Price
	sale := component.SaleData{
		EndBlock:         blockHeight + 1000,
		BuyingAsset:      common.ConstantID,
		BuyingAmount:     burnAmount,
		DefaultBuyPrice:  constantPrice,
		SellingAsset:     bondToSell.BondID,
		SellingAmount:    sellingAmount,
		DefaultSellPrice: bondToSell.Price,
	}
	return []component.SaleData{sale}, nil
}

func buildTradeSellBond(
	burnAmount uint64,
	blockHeight uint64,
	bonds []*entities.DCBBondInfo,
) ([]*component.TradeBondWithGOV, error) {
	// Check if there's a bond that can reduce burnAmount of Constant
	var bondToSell *entities.DCBBondInfo
	for _, b := range bonds {
		if b.Maturity < blockHeight && b.BuyBack*b.Amount >= burnAmount {
			bondToSell = b
			break
		}
	}
	if bondToSell == nil {
		return nil, nil
	}

	amount := 1 + (burnAmount-1)/bondToSell.BuyBack
	trade := &component.TradeBondWithGOV{
		BondID: &bondToSell.BondID,
		Amount: amount,
		Buy:    false,
	}
	return []*component.TradeBondWithGOV{trade}, nil
}

func buildSpendReserve(
	burnAmount uint64,
	constantPrice uint64,
	blockHeight uint64,
	dr *DataRequester,
) (map[common.Hash]*component.SpendReserveData, error) {
	// TODO(@0xbunyip): choose between ETH and USD
	price, err := dr.AssetPrice(common.ETHAssetID)
	if err != nil {
		return nil, err
	}
	reserve := map[common.Hash]*component.SpendReserveData{
		common.ETHAssetID: &component.SpendReserveData{
			EndBlock:        blockHeight + 1000,
			ReserveMinPrice: price,
			Amount:          burnAmount,
		},
	}
	return reserve, nil
}

func buildCrowdsalesBuyBond(
	mintAmount uint64,
	constantPrice uint64,
	blockHeight uint64,
	dr *DataRequester,
) ([]component.SaleData, error) {
	bonds, err := dr.BondsCirculating()
	if err != nil {
		return nil, err
	}

	// Check if there's a bond that can mint mintAmount of Constant
	var bondToBuy *entities.DCBBondInfo
	for _, b := range bonds {
		if b.Price > 0 && b.Price*b.Amount >= mintAmount {
			bondToBuy = b
			break
		}
	}
	if bondToBuy == nil {
		return nil, nil
	}

	buyingAmount := 1 + (mintAmount-1)/bondToBuy.Price
	sale := component.SaleData{
		EndBlock:         blockHeight + 1000,
		BuyingAsset:      bondToBuy.BondID,
		BuyingAmount:     buyingAmount,
		DefaultBuyPrice:  bondToBuy.Price,
		SellingAsset:     common.ConstantID,
		SellingAmount:    mintAmount,
		DefaultSellPrice: constantPrice,
	}
	return []component.SaleData{sale}, nil
}

func buildTradeBuyBond(
	mintAmount uint64,
	blockHeight uint64,
	dr *DataRequester,
) ([]*component.TradeBondWithGOV, error) {
	// Check if GOV's selling bond can cover mintAmount of Constant
	bondToBuy, err := dr.CurrentSellingBond()
	if err != nil || bondToBuy == nil {
		return nil, err
	}

	amount := 1 + (mintAmount-1)/bondToBuy.Price
	trade := &component.TradeBondWithGOV{
		BondID: &bondToBuy.BondID,
		Amount: amount,
		Buy:    true,
	}
	return []*component.TradeBondWithGOV{trade}, nil
}
