/**
The MIT License (MIT)

Copyright (c) 2018 SmartContract ChainLink, Ltd.
*/
// Package utils is used for common functions and tools used across the codebase.
package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/go-wallet-sdk/util/abi"
	"golang.org/x/crypto/sha3"
	"math/big"
)

const (
	// DefaultSecretSize is the entropy in bytes to generate a base64 string of 64 characters.
	DefaultSecretSize = 48
	// EVMWordByteLen the length of an EVM Word Byte
	EVMWordByteLen = 32
)

// Keccak256 is a simplified interface for the legacy SHA3 implementation that
// Ethereum uses.
func Keccak256(in []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(in)
	return hash.Sum(nil), err
}

// MustHash returns the keccak256 hash, or panics on failure.
func MustHash(in string) common.Hash {
	out, err := Keccak256([]byte(in))
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(out)
}

// HexToBig parses the given hex string or panics if it is invalid.
func HexToBig(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(fmt.Errorf(`failed to convert "%s" as hex to big.Int`, s))
	}
	return n
}

// Uint256ToBytes32 returns the bytes32 encoding of the big int provided
func Uint256ToBytes32(n *big.Int) []byte {
	if n.BitLen() > 256 {
		panic("vrf.uint256ToBytes32: too big to marshal to uint256")
	}
	return common.LeftPadBytes(n.Bytes(), 32)
}

// Uint256ToBytes is x represented as the bytes of a uint256
func Uint256ToBytes(x *big.Int) (uint256 []byte, err error) {
	if x.Cmp(abi.MaxBig256) > 0 {
		return nil, fmt.Errorf("too large to convert to uint256")
	}
	uint256 = common.LeftPadBytes(x.Bytes(), EVMWordByteLen)
	if x.Cmp(big.NewInt(0).SetBytes(uint256)) != 0 {
		panic("failed to round-trip uint256 back to source big.Int")
	}
	return uint256, err
}
