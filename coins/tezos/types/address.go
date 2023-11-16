/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// AddressType represents the type of a Tezos signature.
type AddressType byte

const (
	AddressTypeInvalid AddressType = iota
	AddressTypeEd25519
	AddressTypeSecp256k1
	AddressTypeP256
	AddressTypeContract
	AddressTypeBlinded
	AddressTypeSapling
	AddressTypeBls12_381
	AddressTypeToru
)

var (
	// ErrChecksumMismatch describes an error where decoding failed due
	// to a bad checksum.
	ErrChecksumMismatch = errors.New("tezos: checksum mismatch")
	// InvalidAddress is an empty invalid address
	InvalidAddress = Address{Type: AddressTypeInvalid, Hash: nil}
)

func (t AddressType) Tag() byte {
	switch t {
	case AddressTypeEd25519:
		return 0
	case AddressTypeSecp256k1:
		return 1
	case AddressTypeP256:
		return 2
	case AddressTypeBlinded:
		return 3
	case AddressTypeBls12_381:
		return 4
	default:
		return 255
	}
}

func (t AddressType) HashType() HashType {
	switch t {
	case AddressTypeEd25519:
		return HashTypePkhEd25519
	case AddressTypeSecp256k1:
		return HashTypePkhSecp256k1
	case AddressTypeP256:
		return HashTypePkhP256
	case AddressTypeContract:
		return HashTypePkhNocurve
	case AddressTypeBlinded:
		return HashTypePkhBlinded
	case AddressTypeSapling:
		return HashTypeSaplingAddress
	case AddressTypeBls12_381:
		return HashTypePkhBls12_381
	case AddressTypeToru:
		return HashTypeToruAddress
	default:
		return HashTypeInvalid
	}
}

type Address struct {
	Type AddressType
	Hash []byte
}

func (a Address) IsValid() bool {
	return a.Type != AddressTypeInvalid && len(a.Hash) == a.Type.HashType().Len()
}

func (a Address) Equal(b Address) bool {
	return a.Type == b.Type && bytes.Equal(a.Hash, b.Hash)
}

func (a Address) String() string {
	s, _ := EncodeAddress(a.Type, a.Hash)
	return s
}

func (a Address) Bytes() []byte {
	switch a.Type {
	case AddressTypeInvalid:
		return nil
	case AddressTypeContract:
		buf := append([]byte{01}, a.Hash...)
		buf = append(buf, byte(0)) // padding
		return buf
	case AddressTypeToru:
		buf := append([]byte{02}, a.Hash...)
		buf = append(buf, byte(0)) // padding
		return buf
	default:
		return append([]byte{a.Type.Tag()}, a.Hash...)
	}
}

// Bytes22 returns the 22 byte tagged and padded binary encoding for contracts
// and EOAs (tz1/2/3). In contrast to Bytes which outputs the 21 byte address for EOAs
// here we add a leading 0-byte.
func (a Address) Bytes22() []byte {
	switch a.Type {
	case AddressTypeInvalid:
		return nil
	case AddressTypeContract:
		buf := append([]byte{01}, a.Hash...)
		buf = append(buf, byte(0)) // padding
		return buf
	case AddressTypeToru:
		buf := append([]byte{02}, a.Hash...)
		buf = append(buf, byte(0)) // padding
		return buf
	default:
		return append([]byte{00, a.Type.Tag()}, a.Hash...)
	}
}

func EncodeAddress(typ AddressType, addrhash []byte) (string, error) {
	if len(addrhash) != 20 {
		return "", fmt.Errorf("tezos: invalid address hash")
	}
	switch typ {
	case AddressTypeEd25519:
		return CheckEncode(addrhash, ED25519_PUBLIC_KEY_HASH_ID), nil
	case AddressTypeSecp256k1:
		return CheckEncode(addrhash, SECP256K1_PUBLIC_KEY_HASH_ID), nil
	case AddressTypeP256:
		return CheckEncode(addrhash, P256_PUBLIC_KEY_HASH_ID), nil
	case AddressTypeContract:
		return CheckEncode(addrhash, NOCURVE_PUBLIC_KEY_HASH_ID), nil
	case AddressTypeBlinded:
		return CheckEncode(addrhash, BLINDED_PUBLIC_KEY_HASH_ID), nil
	case AddressTypeSapling:
		return CheckEncode(addrhash, SAPLING_ADDRESS_ID), nil
	case AddressTypeBls12_381:
		return CheckEncode(addrhash, BLS12_381_PUBLIC_KEY_HASH_ID), nil
	case AddressTypeToru:
		return CheckEncode(addrhash, TORU_ADDRESS_ID), nil
	default:
		return "", fmt.Errorf("tezos: unknown address type %s for hash=%x", typ, addrhash)
	}
}

func ParseAddress(addr string) (Address, error) {
	if len(addr) == 0 {
		return InvalidAddress, nil
	}
	a := Address{}
	sz := 3
	if strings.HasPrefix(addr, BLINDED_PUBLIC_KEY_HASH_PREFIX) ||
		strings.HasPrefix(addr, TORU_ADDRESS_PREFIX) {
		sz = 4
	}
	decoded, version, err := CheckDecode(addr, sz, nil)
	if err != nil {
		if err == ErrChecksum {
			return a, ErrChecksumMismatch
		}
		return a, fmt.Errorf("tezos: decoded address is of unknown format: %w", err)
	}
	if len(decoded) != 20 {
		return a, errors.New("tezos: decoded address hash is of invalid length")
	}
	switch true {
	case bytes.Equal(version, ED25519_PUBLIC_KEY_HASH_ID):
		return Address{Type: AddressTypeEd25519, Hash: decoded}, nil
	case bytes.Equal(version, SECP256K1_PUBLIC_KEY_HASH_ID):
		return Address{Type: AddressTypeSecp256k1, Hash: decoded}, nil
	case bytes.Equal(version, P256_PUBLIC_KEY_HASH_ID):
		return Address{Type: AddressTypeP256, Hash: decoded}, nil
	case bytes.Equal(version, NOCURVE_PUBLIC_KEY_HASH_ID):
		return Address{Type: AddressTypeContract, Hash: decoded}, nil
	case bytes.Equal(version, BLINDED_PUBLIC_KEY_HASH_ID):
		return Address{Type: AddressTypeBlinded, Hash: decoded}, nil
	case bytes.Equal(version, SAPLING_ADDRESS_ID):
		return Address{Type: AddressTypeSapling, Hash: decoded}, nil
	case bytes.Equal(version, BLS12_381_PUBLIC_KEY_HASH_ID):
		return Address{Type: AddressTypeBls12_381, Hash: decoded}, nil
	case bytes.Equal(version, TORU_ADDRESS_ID):
		return Address{Type: AddressTypeToru, Hash: decoded}, nil
	default:
		return a, fmt.Errorf("tezos: decoded address %s is of unknown type %x", addr, version)
	}
}
