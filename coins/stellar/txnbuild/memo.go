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
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/stellar/support/errors"
	"github.com/okx/go-wallet-sdk/coins/stellar/xdr"
)

// MemoText is used to send human messages of up to 28 bytes of ASCII/UTF-8.
type MemoText string

// MemoID is an identifier representing the transaction originator.
type MemoID uint64

// MemoHash is a hash representing a reference to another transaction.
type MemoHash [32]byte

// MemoReturn is a hash representing the hash of the transaction the sender is refunding.
type MemoReturn [32]byte

// MemoTextMaxLength is the maximum number of bytes allowed for a text memo.
const MemoTextMaxLength = 28

// Memo represents the superset of all memo types.
type Memo interface {
	ToXDR() (xdr.Memo, error)
}

// ToXDR for MemoText returns an XDR object representation of a Memo of the same type.
func (mt MemoText) ToXDR() (xdr.Memo, error) {
	if len(mt) > MemoTextMaxLength {
		return xdr.Memo{}, fmt.Errorf("Memo text can't be longer than %d bytes", MemoTextMaxLength)
	}

	return xdr.NewMemo(xdr.MemoTypeMemoText, string(mt))
}

// ToXDR for MemoID returns an XDR object representation of a Memo of the same type.
func (mid MemoID) ToXDR() (xdr.Memo, error) {
	return xdr.NewMemo(xdr.MemoTypeMemoId, xdr.Uint64(mid))
}

// ToXDR for MemoHash returns an XDR object representation of a Memo of the same type.
func (mh MemoHash) ToXDR() (xdr.Memo, error) {
	return xdr.NewMemo(xdr.MemoTypeMemoHash, xdr.Hash(mh))
}

// ToXDR for MemoReturn returns an XDR object representation of a Memo of the same type.
func (mr MemoReturn) ToXDR() (xdr.Memo, error) {
	return xdr.NewMemo(xdr.MemoTypeMemoReturn, xdr.Hash(mr))
}

// memoFromXDR returns a Memo from XDR
func memoFromXDR(memo xdr.Memo) (Memo, error) {
	var newMemo Memo
	var memoCreated bool

	switch memo.Type {
	case xdr.MemoTypeMemoText:
		value, ok := memo.GetText()
		newMemo = MemoText(value)
		memoCreated = ok
	case xdr.MemoTypeMemoId:
		value, ok := memo.GetId()
		newMemo = MemoID(uint64(value))
		memoCreated = ok
	case xdr.MemoTypeMemoHash:
		value, ok := memo.GetHash()
		newMemo = MemoHash(value)
		memoCreated = ok
	case xdr.MemoTypeMemoReturn:
		value, ok := memo.GetRetHash()
		newMemo = MemoReturn(value)
		memoCreated = ok
	case xdr.MemoTypeMemoNone:
		memoCreated = true
	}

	if !memoCreated {
		return nil, errors.New("invalid memo")
	}

	return newMemo, nil
}
