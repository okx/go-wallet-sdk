package txnbuild

import (
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
)

// ChangeTrustAsset represents a Stellar change trust asset.
type ChangeTrustAsset interface {
	BasicAsset
	ToXDR() (xdr.ChangeTrustAsset, error)
	ToChangeTrustAsset() (ChangeTrustAsset, error)
	ToTrustLineAsset() (TrustLineAsset, error)
}

// ChangeTrustAssetWrapper wraps a native/credit Asset so it generates xdr to be used in a change trust operation.
type ChangeTrustAssetWrapper struct {
	Asset
}

// ToXDR for ChangeTrustAssetWrapper generates the xdr.TrustLineAsset.
func (ctaw ChangeTrustAssetWrapper) ToXDR() (xdr.ChangeTrustAsset, error) {
	x, err := ctaw.Asset.ToXDR()
	if err != nil {
		return xdr.ChangeTrustAsset{}, err
	}
	return x.ToChangeTrustAsset(), nil
}

func assetFromChangeTrustAssetXDR(xAsset xdr.ChangeTrustAsset) (ChangeTrustAsset, error) {
	a, err := assetFromXDR(xAsset.ToAsset())
	if err != nil {
		return nil, err
	}
	return ChangeTrustAssetWrapper{a}, nil
}
