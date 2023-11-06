package atomic_market

import (
	"github.com/okx/go-wallet-sdk/coins/eos/types"
)

func NewAssertSale(from, listingPriceToAssert, settlementSymbolToAssert string,
	saleId uint64, assetIdsToAssert []uint64) *types.Action {
	// TODO: this is only for WAX
	asset, _ := types.NewAssetFromString(listingPriceToAssert, types.WAXSymbol)
	return &types.Action{
		Account: "atomicmarket",
		Name:    "assertsale",
		Authorization: []types.PermissionLevel{
			{Actor: types.AN(from), Permission: "active"},
		},
		ActionData: types.NewActionData(AssertSale{
			SaleId:                   saleId,
			ListingPriceToAssert:     asset,
			SettlementSymbolToAssert: types.WAXSymbol,
			AssetIdsToAssert:         assetIdsToAssert,
		}),
	}
}

// AssertSale represents the `assertsale` action.
type AssertSale struct {
	SaleId                   uint64       `json:"sale_id"`
	AssetIdsToAssert         []uint64     `json:"asset_ids_to_assert"`
	ListingPriceToAssert     types.Asset  `json:"listing_price_to_assert"`
	SettlementSymbolToAssert types.Symbol `json:"settlement_symbol_to_assert"`
}

func NewPurchaseSale(from string, saleId, intendedDelphiMedian uint64, takerMarketplace string) *types.Action {
	return &types.Action{
		Account: "atomicmarket",
		Name:    "purchasesales",
		Authorization: []types.PermissionLevel{
			{Actor: types.AN(from), Permission: "active"},
		},
		ActionData: types.NewActionData(PurchaseSale{
			SaleId:               saleId,
			IntendedDelphiMedian: intendedDelphiMedian,
			TakerMarketplace:     types.AN(takerMarketplace),
			Buyer:                types.AN(from),
		}),
	}
}

type PurchaseSale struct {
	Buyer                types.AccountName `json:"buyer"`
	SaleId               uint64            `json:"sale_id"`
	IntendedDelphiMedian uint64            `json:"intended_delphi_median"`
	TakerMarketplace     types.AccountName `json:"taker_marketplace"`
}
