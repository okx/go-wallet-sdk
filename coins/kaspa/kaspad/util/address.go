// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package util

import (
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/util/bech32"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

var (
	// ErrUnknownAddressType describes an error where an address can not
	// decoded as a specific address type due to the string encoding
	// begining with an identifier byte unknown to any standard or
	// registered (via dagconfig.Register) network.
	ErrUnknownAddressType = errors.New("unknown address type")
)

const (
	// PubKey addresses always have the version byte set to 0.
	pubKeyAddrID = 0x00

	// PubKey addresses always have the version byte set to 1.
	pubKeyECDSAAddrID = 0x01

	// ScriptHash addresses always have the version byte set to 8.
	scriptHashAddrID = 0x08
)

// Bech32Prefix is the human-readable prefix for a Bech32 address.
type Bech32Prefix int

// Constants that define Bech32 address prefixes. Every network is assigned
// a unique prefix.
const (
	// Unknown/Erroneous prefix
	Bech32PrefixUnknown Bech32Prefix = iota

	// Prefix for the main network.
	Bech32PrefixKaspa

	// Prefix for the dev network.
	Bech32PrefixKaspaDev

	// Prefix for the test network.
	Bech32PrefixKaspaTest

	// Prefix for the simulation network.
	Bech32PrefixKaspaSim
)

// Map from strings to Bech32 address prefix constants for parsing purposes.
var stringsToBech32Prefixes = map[string]Bech32Prefix{
	"kaspa":     Bech32PrefixKaspa,
	"kaspadev":  Bech32PrefixKaspaDev,
	"kaspatest": Bech32PrefixKaspaTest,
	"kaspasim":  Bech32PrefixKaspaSim,
}

// ParsePrefix attempts to parse a Bech32 address prefix.
func ParsePrefix(prefixString string) (Bech32Prefix, error) {
	prefix, ok := stringsToBech32Prefixes[prefixString]
	if !ok {
		return Bech32PrefixUnknown, errors.Errorf("could not parse prefix %s", prefixString)
	}

	return prefix, nil
}

// Converts from Bech32 address prefixes to their string values
func (prefix Bech32Prefix) String() string {
	for key, value := range stringsToBech32Prefixes {
		if prefix == value {
			return key
		}
	}

	return ""
}

// encodeAddress returns a human-readable payment address given a network prefix
// and a payload which encodes the kaspa network and address type. It is used
// in both pay-to-pubkey (P2PK) and pay-to-script-hash (P2SH) address
// encoding.
func encodeAddress(prefix Bech32Prefix, payload []byte, version byte) string {
	return bech32.Encode(prefix.String(), payload, version)
}

// Address is an interface type for any type of destination a transaction
// output may spend to. This includes pay-to-pubkey (P2PK)
// and pay-to-script-hash (P2SH). Address is designed to be generic
// enough that other kinds of addresses may be added in the future without
// changing the decoding and encoding API.
type Address interface {
	// String returns the string encoding of the transaction output
	// destination.
	//
	// Please note that String differs subtly from EncodeAddress: String
	// will return the value as a string without any conversion, while
	// EncodeAddress may convert destination types (for example,
	// converting pubkeys to P2PK addresses) before encoding as a
	// payment address string.
	String() string

	// EncodeAddress returns the string encoding of the payment address
	// associated with the Address value. See the comment on String
	// for how this method differs from String.
	EncodeAddress() string

	// ScriptAddress returns the raw bytes of the address to be used
	// when inserting the address into a txout's script.
	ScriptAddress() []byte

	// Prefix returns the prefix for this address
	Prefix() Bech32Prefix

	// IsForPrefix returns whether or not the address is associated with the
	// passed kaspa network.
	IsForPrefix(prefix Bech32Prefix) bool
}

// DecodeAddress decodes the string encoding of an address and returns
// the Address if addr is a valid encoding for a known address type.
//
// If any expectedPrefix except Bech32PrefixUnknown is passed, it is compared to the
// prefix extracted from the address, and if the two do not match - an error is returned
func DecodeAddress(addr string, expectedPrefix Bech32Prefix) (Address, error) {
	prefixString, decoded, version, err := bech32.Decode(addr)
	if err != nil {
		return nil, errors.Errorf("decoded address is of unknown format: %s", err)
	}

	prefix, err := ParsePrefix(prefixString)
	if err != nil {
		return nil, errors.Errorf("decoded address's prefix could not be parsed: %s", err)
	}
	if expectedPrefix != Bech32PrefixUnknown && expectedPrefix != prefix {
		return nil, errors.Errorf("decoded address is of wrong network. Expected %s but got %s", expectedPrefix,
			prefix)
	}

	switch version {
	case pubKeyAddrID:
		return newAddressPubKey(prefix, decoded)
	case pubKeyECDSAAddrID:
		return newAddressPubKeyECDSA(prefix, decoded)
	case scriptHashAddrID:
		return newAddressScriptHashFromHash(prefix, decoded)
	default:
		return nil, ErrUnknownAddressType
	}
}

// PublicKeySize is the public key size for a schnorr public key
const PublicKeySize = 32

// AddressPublicKey is an Address for a pay-to-pubkey (P2PK)
// transaction.
type AddressPublicKey struct {
	prefix    Bech32Prefix
	publicKey [PublicKeySize]byte
}

// NewAddressPublicKey returns a new AddressPublicKey. publicKey must be 32
// bytes.
func NewAddressPublicKey(publicKey []byte, prefix Bech32Prefix) (*AddressPublicKey, error) {
	return newAddressPubKey(prefix, publicKey)
}

// newAddressPubKey is the internal API to create a pubkey address
// with a known leading identifier byte for a network, rather than looking
// it up through its parameters. This is useful when creating a new address
// structure from a string encoding where the identifier byte is already
// known.
func newAddressPubKey(prefix Bech32Prefix, publicKey []byte) (*AddressPublicKey, error) {
	// Check for a valid pubkey length.
	if len(publicKey) != PublicKeySize {
		return nil, errors.Errorf("publicKey must be %d bytes", PublicKeySize)
	}

	addr := &AddressPublicKey{prefix: prefix}
	copy(addr.publicKey[:], publicKey)
	return addr, nil
}

// EncodeAddress returns the string encoding of a pay-to-pubkey
// address. Part of the Address interface.
func (a *AddressPublicKey) EncodeAddress() string {
	return encodeAddress(a.prefix, a.publicKey[:], pubKeyAddrID)
}

// ScriptAddress returns the bytes to be included in a txout script to pay
// to a pubkey. Part of the Address interface.
func (a *AddressPublicKey) ScriptAddress() []byte {
	return a.publicKey[:]
}

// IsForPrefix returns whether or not the pay-to-pubkey address is associated
// with the passed kaspa network.
func (a *AddressPublicKey) IsForPrefix(prefix Bech32Prefix) bool {
	return a.prefix == prefix
}

// Prefix returns the prefix for this address
func (a *AddressPublicKey) Prefix() Bech32Prefix {
	return a.prefix
}

// String returns a human-readable string for the pay-to-pubkey address.
// This is equivalent to calling EncodeAddress, but is provided so the type can
// be used as a fmt.Stringer.
func (a *AddressPublicKey) String() string {
	return a.EncodeAddress()
}

// PublicKeySizeECDSA is the public key size for an ECDSA public key
const PublicKeySizeECDSA = 33

// AddressPublicKeyECDSA is an Address for a pay-to-pubkey (P2PK)
// ECDSA transaction.
type AddressPublicKeyECDSA struct {
	prefix    Bech32Prefix
	publicKey [PublicKeySizeECDSA]byte
}

// NewAddressPublicKeyECDSA returns a new AddressPublicKeyECDSA. publicKey must be 33
// bytes.
func NewAddressPublicKeyECDSA(publicKey []byte, prefix Bech32Prefix) (*AddressPublicKeyECDSA, error) {
	return newAddressPubKeyECDSA(prefix, publicKey)
}

// newAddressPubKeyECDSA is the internal API to create an ECDSA pubkey address
// with a known leading identifier byte for a network, rather than looking
// it up through its parameters. This is useful when creating a new address
// structure from a string encoding where the identifier byte is already
// known.
func newAddressPubKeyECDSA(prefix Bech32Prefix, publicKey []byte) (*AddressPublicKeyECDSA, error) {
	// Check for a valid pubkey length.
	if len(publicKey) != PublicKeySizeECDSA {
		return nil, errors.Errorf("publicKey must be %d bytes", PublicKeySizeECDSA)
	}

	addr := &AddressPublicKeyECDSA{prefix: prefix}
	copy(addr.publicKey[:], publicKey)
	return addr, nil
}

// EncodeAddress returns the string encoding of a pay-to-pubkey
// address. Part of the Address interface.
func (a *AddressPublicKeyECDSA) EncodeAddress() string {
	return encodeAddress(a.prefix, a.publicKey[:], pubKeyECDSAAddrID)
}

// ScriptAddress returns the bytes to be included in a txout script to pay
// to a pubkey. Part of the Address interface.
func (a *AddressPublicKeyECDSA) ScriptAddress() []byte {
	return a.publicKey[:]
}

// IsForPrefix returns whether or not the pay-to-pubkey address is associated
// with the passed kaspa network.
func (a *AddressPublicKeyECDSA) IsForPrefix(prefix Bech32Prefix) bool {
	return a.prefix == prefix
}

// Prefix returns the prefix for this address
func (a *AddressPublicKeyECDSA) Prefix() Bech32Prefix {
	return a.prefix
}

// String returns a human-readable string for the pay-to-pubkey address.
// This is equivalent to calling EncodeAddress, but is provided so the type can
// be used as a fmt.Stringer.
func (a *AddressPublicKeyECDSA) String() string {
	return a.EncodeAddress()
}

// AddressScriptHash is an Address for a pay-to-script-publicKey (P2SH)
// transaction.
type AddressScriptHash struct {
	prefix Bech32Prefix
	hash   [blake2b.Size256]byte
}

// NewAddressScriptHash returns a new AddressScriptHash.
func NewAddressScriptHash(serializedScript []byte, prefix Bech32Prefix) (*AddressScriptHash, error) {
	scriptHash := HashBlake2b(serializedScript)
	return newAddressScriptHashFromHash(prefix, scriptHash)
}

// NewAddressScriptHashFromHash returns a new AddressScriptHash. scriptHash
// must be 20 bytes.
func NewAddressScriptHashFromHash(scriptHash []byte, prefix Bech32Prefix) (*AddressScriptHash, error) {
	return newAddressScriptHashFromHash(prefix, scriptHash)
}

// newAddressScriptHashFromHash is the internal API to create a script hash
// address with a known leading identifier byte for a network, rather than
// looking it up through its parameters. This is useful when creating a new
// address structure from a string encoding where the identifer byte is already
// known.
func newAddressScriptHashFromHash(prefix Bech32Prefix, scriptHash []byte) (*AddressScriptHash, error) {
	// Check for a valid script hash length.
	if len(scriptHash) != blake2b.Size256 {
		return nil, errors.Errorf("scriptHash must be %d bytes", blake2b.Size256)
	}

	addr := &AddressScriptHash{prefix: prefix}
	copy(addr.hash[:], scriptHash)
	return addr, nil
}

// EncodeAddress returns the string encoding of a pay-to-script-hash
// address. Part of the Address interface.
func (a *AddressScriptHash) EncodeAddress() string {
	return encodeAddress(a.prefix, a.hash[:], scriptHashAddrID)
}

// ScriptAddress returns the bytes to be included in a txout script to pay
// to a script hash. Part of the Address interface.
func (a *AddressScriptHash) ScriptAddress() []byte {
	return a.hash[:]
}

// IsForPrefix returns whether or not the pay-to-script-hash address is associated
// with the passed kaspa network.
func (a *AddressScriptHash) IsForPrefix(prefix Bech32Prefix) bool {
	return a.prefix == prefix
}

// Prefix returns the prefix for this address
func (a *AddressScriptHash) Prefix() Bech32Prefix {
	return a.prefix
}

// String returns a human-readable string for the pay-to-script-hash address.
// This is equivalent to calling EncodeAddress, but is provided so the type can
// be used as a fmt.Stringer.
func (a *AddressScriptHash) String() string {
	return a.EncodeAddress()
}

// HashBlake2b returns the underlying array of the script hash. This can be useful
// when an array is more appropiate than a slice (for example, when used as map
// keys).
func (a *AddressScriptHash) HashBlake2b() *[blake2b.Size256]byte {
	return &a.hash
}
