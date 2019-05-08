package agents

const (
	// DefaultFee - default fee
	DefaultFee = 20

	// DefaultExecuteDuration - default execute duration
	DefaultExecuteDuration = 100

	// Pegging price of Constant
	Peg = 100

	// Acceptable range of price to be above/below the peg
	PriceCeiling = 10
	PriceFloor   = 10

	// Min amount of Constant needed to burn/mint in order to submit a proposal
	MinAmountToBurn = uint64(1000)
	MinAmountToMint = uint64(1000)

	// Maximum number of terms with crowdsales and trades to try before moving
	// on to the next tool
	NumSaleToTry  = 3
	NumTradeToTry = 2
)
