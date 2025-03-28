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
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/stellar/strkey"
)

// Address returns the strkey encoded form of this AccountId.  This method will
// panic if the accountid is backed by a public key of an unknown type.
func (aid AccountId) Address() string {
	address, err := aid.GetAddress()
	if err != nil {
		panic(err)
	}
	return address
}

// GetAddress returns the strkey encoded form of this AccountId, and an error
// if the AccountId is backed by a public key of an unknown type.
func (aid *AccountId) GetAddress() (string, error) {
	if aid == nil {
		return "", nil
	}

	switch aid.Type {
	case PublicKeyTypePublicKeyTypeEd25519:
		ed, ok := aid.GetEd25519()
		if !ok {
			return "", fmt.Errorf("Could not get Ed25519")
		}
		raw := make([]byte, 32)
		copy(raw, ed[:])
		return strkey.Encode(strkey.VersionByteAccountID, raw)
	default:
		return "", fmt.Errorf("Unknown account id type: %v", aid.Type)
	}
}

// Equals returns true if `other` is equivalent to `aid`
func (aid *AccountId) Equals(other AccountId) bool {
	if aid.Type != other.Type {
		return false
	}

	switch aid.Type {
	case PublicKeyTypePublicKeyTypeEd25519:
		l := aid.MustEd25519()
		r := other.MustEd25519()
		return l == r
	default:
		panic(fmt.Errorf("Unknown account id type: %v", aid.Type))
	}
}

// LedgerKey implements the `Keyer` interface
func (aid *AccountId) LedgerKey() (key LedgerKey, err error) {
	return key, nil
}

func (e *EncodingBuffer) accountIdCompressEncodeTo(aid AccountId) error {
	if err := e.xdrEncoderBuf.WriteByte(byte(aid.Type)); err != nil {
		return err
	}
	switch aid.Type {
	case PublicKeyTypePublicKeyTypeEd25519:
		_, err := e.xdrEncoderBuf.Write(aid.Ed25519[:])
		return err
	default:
		panic("Unknown type")
	}
}

func MustAddress(address string) AccountId {
	aid := AccountId{}
	err := aid.SetAddress(address)
	if err != nil {
		panic(err)
	}
	return aid
}

func MustAddressPtr(address string) *AccountId {
	aid := MustAddress(address)
	return &aid
}

// AddressToAccountId returns an AccountId for a given address string.
// If the address is not valid the error returned will not be nil
func AddressToAccountId(address string) (AccountId, error) {
	result := AccountId{}
	err := result.SetAddress(address)

	return result, err
}

// SetAddress modifies the receiver, setting it's value to the AccountId form
// of the provided address.
func (aid *AccountId) SetAddress(address string) error {
	if aid == nil {
		return nil
	}

	raw, err := strkey.Decode(strkey.VersionByteAccountID, address)
	if err != nil {
		return err
	}

	if len(raw) != 32 {
		return errors.New("invalid address")
	}

	var ui Uint256
	copy(ui[:], raw)

	*aid, err = NewAccountId(PublicKeyTypePublicKeyTypeEd25519, ui)

	return err
}

// ToMuxedAccount transforms an AccountId into a MuxedAccount with
// a zero memo id
func (aid *AccountId) ToMuxedAccount() MuxedAccount {
	result := MuxedAccount{Type: CryptoKeyTypeKeyTypeEd25519}
	switch aid.Type {
	case PublicKeyTypePublicKeyTypeEd25519:
		result.Ed25519 = aid.Ed25519
	default:
		panic(fmt.Errorf("Unknown account id type: %v", aid.Type))
	}
	return result
}
