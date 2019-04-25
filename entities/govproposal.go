package entities

type GOVProposalRes struct {
	RPCBaseRes
	Result interface{}
}

type GOVProposal struct {
	GOVParams         *GOVParams
	ExecuteDuration   uint64
	Explanation       string
	PaymentAddress    string
	ConstitutionIndex uint64
}

type GOVParams struct {
	SalaryPerTx      uint64 // salary for each tx in block(mili constant)
	BasicSalary      uint64 // basic salary per block(mili constant)
	FeePerKbTx       uint64
	SellingBonds     *SellingBonds
	SellingGOVTokens *SellingGOVTokens
	RefundInfo       *RefundInfo
	OracleNetwork    *OracleNetwork
}

type SellingBonds struct {
	BondName       string
	BondSymbol     string
	TotalIssue     uint64
	BondsToSell    uint64
	BondPrice      uint64 // in Constant unit
	Maturity       uint64
	BuyBackPrice   uint64 // in Constant unit
	StartSellingAt uint64 // start selling bonds at block height
	SellingWithin  uint64 // selling bonds within n blocks
}

type SellingGOVTokens struct {
	TotalIssue      uint64
	GOVTokensToSell uint64
	GOVTokenPrice   uint64 // in Constant unit
	StartSellingAt  uint64 // start selling gov tokens at block height
	SellingWithin   uint64 // selling tokens within n blocks
}

type OracleNetwork struct {
	OraclePubKeys          []string // hex string encoded
	WrongTimesAllowed      uint8
	Quorum                 uint8
	AcceptableErrorMargin  uint32
	UpdateFrequency        uint32
	OracleRewardMultiplier uint8
}

type RefundInfo struct {
	ThresholdToLargeTx uint64
	RefundAmount       uint64
}
