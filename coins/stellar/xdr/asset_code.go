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
