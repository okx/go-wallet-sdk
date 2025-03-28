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

// ToAsset converts ChangeTrustAsset to Asset. Panics on type other than
// AssetTypeAssetTypeNative, AssetTypeAssetTypeCreditAlphanum4 or
// AssetTypeAssetTypeCreditAlphanum12.
func (tla ChangeTrustAsset) ToAsset() Asset {
	var a Asset

	a.Type = tla.Type

	switch a.Type {
	case AssetTypeAssetTypeNative:
		// Empty branch
	case AssetTypeAssetTypeCreditAlphanum4:
		assetCode4 := *tla.AlphaNum4
		a.AlphaNum4 = &assetCode4
	case AssetTypeAssetTypeCreditAlphanum12:
		assetCode12 := *tla.AlphaNum12
		a.AlphaNum12 = &assetCode12
	default:
		panic(fmt.Errorf("Cannot transform type %v to Asset", a.Type))
	}

	return a
}
