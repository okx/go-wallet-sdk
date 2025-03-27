package txnbuild

import (
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
)

// TrustLineAsset represents a Stellar trust line asset.
type TrustLineAsset interface {
	BasicAsset
	ToXDR() (xdr.TrustLineAsset, error)
}

// TrustLineAssetWrapper wraps a native/credit Asset so it generates xdr to be used in a trust line operation.
type TrustLineAssetWrapper struct {
	Asset
}

// ToXDR for TrustLineAssetWrapper generates the xdr.TrustLineAsset.
func (tlaw TrustLineAssetWrapper) ToXDR() (xdr.TrustLineAsset, error) {
	x, err := tlaw.Asset.ToXDR()
	if err != nil {
		return xdr.TrustLineAsset{}, err
	}
	return x.ToTrustLineAsset(), nil
}
