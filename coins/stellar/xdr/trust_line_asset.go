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

// ToAsset converts TrustLineAsset to Asset. Panics on type other than
// AssetTypeAssetTypeNative, AssetTypeAssetTypeCreditAlphanum4 or
// AssetTypeAssetTypeCreditAlphanum12.
func (tla TrustLineAsset) ToAsset() Asset {
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

// MustExtract behaves as Extract, but panics if an error occurs.
func (a TrustLineAsset) Extract(typ interface{}, code interface{}, issuer interface{}) error {
	return a.ToAsset().Extract(typ, code, issuer)
}

// MustExtract behaves as Extract, but panics if an error occurs.
func (a TrustLineAsset) MustExtract(typ interface{}, code interface{}, issuer interface{}) {
	err := a.ToAsset().Extract(typ, code, issuer)

	if err != nil {
		panic(err)
	}
}

func trimRightZeros(b []byte) []byte {
	if len(b) == 0 {
		return b
	}
	i := len(b)
	for ; i > 0; i-- {
		if b[i-1] != 0 {
			break
		}
	}
	return b[:i]
}

func (e *EncodingBuffer) assetTrustlineCompressEncodeTo(a TrustLineAsset) error {
	if err := e.xdrEncoderBuf.WriteByte(byte(a.Type)); err != nil {
		return err
	}

	switch a.Type {
	case AssetTypeAssetTypeNative:
		return nil
	case AssetTypeAssetTypeCreditAlphanum4:
		code := trimRightZeros(a.AlphaNum4.AssetCode[:])
		if _, err := e.xdrEncoderBuf.Write(code); err != nil {
			return err
		}
		return e.accountIdCompressEncodeTo(a.AlphaNum4.Issuer)
	case AssetTypeAssetTypeCreditAlphanum12:
		code := trimRightZeros(a.AlphaNum12.AssetCode[:])
		if _, err := e.xdrEncoderBuf.Write(code); err != nil {
			return err
		}
		return e.accountIdCompressEncodeTo(a.AlphaNum12.Issuer)
	case AssetTypeAssetTypePoolShare:
		_, err := e.xdrEncoderBuf.Write(a.LiquidityPoolId[:])
		return err
	default:
		panic(fmt.Errorf("Unknown asset type: %v", a.Type))
	}
}

func (a TrustLineAsset) Equals(other TrustLineAsset) bool {
	if a.Type != other.Type {
		return false
	}
	switch a.Type {
	case AssetTypeAssetTypeNative,
		AssetTypeAssetTypeCreditAlphanum4,
		AssetTypeAssetTypeCreditAlphanum12:
		// Safe because a.Type == other.Type
		return a.ToAsset().Equals(other.ToAsset())
	case AssetTypeAssetTypePoolShare:
		return *a.LiquidityPoolId == *other.LiquidityPoolId
	default:
		panic(fmt.Errorf("Unknown asset type: %v", a.Type))
	}
}
