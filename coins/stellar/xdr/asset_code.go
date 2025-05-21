package xdr

import (
	"fmt"
)

// ToAsset for AssetCode converts the xdr.AssetCode to a standard xdr.Asset.
func (a AssetCode) ToAsset(issuer AccountId) (asset Asset) {
	var err error

	switch a.Type {
	case AssetTypeAssetTypeCreditAlphanum4:
		asset, err = NewAsset(AssetTypeAssetTypeCreditAlphanum4, AlphaNum4{
			AssetCode: a.MustAssetCode4(),
			Issuer:    issuer,
		})
	case AssetTypeAssetTypeCreditAlphanum12:
		asset, err = NewAsset(AssetTypeAssetTypeCreditAlphanum12, AlphaNum12{
			AssetCode: a.MustAssetCode12(),
			Issuer:    issuer,
		})
	default:
		err = fmt.Errorf("Unexpected type for AllowTrustOpAsset: %d", a.Type)
	}

	if err != nil {
		panic(err)
	}
	return
}
