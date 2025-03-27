package txnbuild

// AssetAmount is a "tuple", pairing an asset with an amount. Used for
// LiquidityPoolDeposit and LiquidityPoolWithdraw operations.
type AssetAmount struct {
	Asset  Asset
	Amount string
}
