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
