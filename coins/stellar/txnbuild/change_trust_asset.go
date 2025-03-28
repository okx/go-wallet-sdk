/*
 * Copyright 2016 Stellar Development Foundation and contributors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file includes portions of third-party code from [https://github.com/stellar/go].
 * The original code is licensed under the Apache License 2.0.
 */

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
