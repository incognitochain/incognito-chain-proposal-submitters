package entities

type CurrentSellingGOVTokens struct {
	GOVTokenID     string `json:"govTokenID"`
	StartSellingAt uint64 `json:"startSellingAt"`
	EndSellingAt   uint64 `json:"endSellingAt"`
	BuyPrice       uint64 `json:"buyPrice"`
	TotalIssue     uint64 `json:"totalIssue"`
	Available      uint64 `json:"available"`
}

type RPCCurrentSellingGOVTokensRes struct {
	RPCBaseRes
	Result *CurrentSellingGOVTokens
}
