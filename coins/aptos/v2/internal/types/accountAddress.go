package types

import (
	"encoding/json"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
)

// AccountAddress a 32-byte representation of an on-chain address
//
// Implements:
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [json.Marshaler]
//   - [json.Unmarshaler]
type AccountAddress [32]byte

// AccountZero is [AccountAddress] 0x0
var AccountZero = AccountAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// AccountOne is [AccountAddress] 0x1
var AccountOne = AccountAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

// AccountTwo is [AccountAddress] 0x2
var AccountTwo = AccountAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

// AccountThree is [AccountAddress] 0x3
var AccountThree = AccountAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}

// AccountFour is [AccountAddress] 0x4
var AccountFour = AccountAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}

// AccountTen is [AccountAddress] 0xA
var AccountTen = AccountAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x0A}

// IsSpecial Returns whether the address is a "special" address. Addresses are considered
// special if the first 63 characters of the hex string are zero. In other words,
// an address is special if the first 31 bytes are zero and the last byte is
// smaller than `0b10000` (16). In other words, special is defined as an address
// that matches the following regex: `^0x0{63}[0-9a-f]$`. In short form this means
// the addresses in the range from `0x0` to `0xf` (inclusive) are special.
// For more details see the v1 address standard defined as part of AIP-40:
// https://github.com/aptos-foundation/AIPs/blob/main/aips/aip-40.md
func (aa *AccountAddress) IsSpecial() bool {
	for _, b := range aa[:31] {
		if b != 0 {
			return false
		}
	}
	return aa[31] < 0x10
}

// String Returns the canonical string representation of the [AccountAddress]
//
// Please use [AccountAddress.StringLong] for all indexer queries.
func (aa *AccountAddress) String() string {
	if aa.IsSpecial() {
		return fmt.Sprintf("0x%x", aa[31])
	} else {
		return util.BytesToHex(aa[:])
	}
}

// FromAuthKey converts [crypto.AuthenticationKey] to [AccountAddress]
func (aa *AccountAddress) FromAuthKey(authKey *crypto.AuthenticationKey) {
	copy(aa[:], authKey[:])
}

// AuthKey converts [AccountAddress] to [crypto.AuthenticationKey]
func (aa *AccountAddress) AuthKey() *crypto.AuthenticationKey {
	authKey := &crypto.AuthenticationKey{}
	copy(authKey[:], aa[:])
	return authKey
}

// StringLong Returns the long string representation of the AccountAddress
//
// This is most commonly used for all indexer queries.
func (aa *AccountAddress) StringLong() string {
	return util.BytesToHex(aa[:])
}

// MarshalBCS Converts the AccountAddress to BCS encoded bytes
func (aa *AccountAddress) MarshalBCS(ser *bcs.Serializer) {
	ser.FixedBytes(aa[:])
}

// UnmarshalBCS Converts the AccountAddress from BCS encoded bytes
func (aa *AccountAddress) UnmarshalBCS(des *bcs.Deserializer) {
	des.ReadFixedBytesInto((*aa)[:])
}

// MarshalJSON converts the AccountAddress to JSON
func (aa *AccountAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// UnmarshalJSON converts the AccountAddress from JSON
func (aa *AccountAddress) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return fmt.Errorf("failed to convert input to AccountAdddress: %w", err)
	}
	err = aa.ParseStringRelaxed(str)
	if err != nil {
		return fmt.Errorf("failed to convert input to AccountAdddress: %w", err)
	}
	return nil
}

// NamedObjectAddress derives a named object address based on the input address as the creator
func (aa *AccountAddress) NamedObjectAddress(seed []byte) (accountAddress AccountAddress) {
	return aa.DerivedAddress(seed, crypto.NamedObjectScheme)
}

// ObjectAddressFromObject derives an object address based on the input address as the creator object
func (aa *AccountAddress) ObjectAddressFromObject(objectAddress *AccountAddress) (accountAddress AccountAddress) {
	return aa.DerivedAddress(objectAddress[:], crypto.DeriveObjectScheme)
}

// ResourceAccount derives an object address based on the input address as the creator
func (aa *AccountAddress) ResourceAccount(seed []byte) (accountAddress AccountAddress) {
	return aa.DerivedAddress(seed, crypto.ResourceAccountScheme)
}

// DerivedAddress addresses are derived by the address, the seed, then the type byte
func (aa *AccountAddress) DerivedAddress(seed []byte, typeByte uint8) (accountAddress AccountAddress) {
	authKey := aa.AuthKey()
	authKey.FromBytesAndScheme(append(authKey[:], seed[:]...), typeByte)
	copy(accountAddress[:], authKey[:])
	return
}
