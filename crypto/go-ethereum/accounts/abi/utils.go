// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/crypto/go-ethereum/common"
	"golang.org/x/crypto/sha3"
	"hash"
	"sync"
)

// ResolveNameConflict returns the next available name for a given thing.
// This helper can be used for lots of purposes:
//
//   - In solidity function overloading is supported, this function can fix
//     the name conflicts of overloaded functions.
//   - In golang binding generation, the parameter(in function, event, error,
//     and struct definition) name will be converted to camelcase style which
//     may eventually lead to name conflicts.
//
// Name conflicts are mostly resolved by adding number suffix. e.g. if the abi contains
// Methods "send" and "send1", ResolveNameConflict would return "send2" for input "send".
func ResolveNameConflict(rawName string, used func(string) bool) string {
	name := rawName
	ok := used(name)
	for idx := 0; ok; idx++ {
		name = fmt.Sprintf("%s%d", rawName, idx)
		ok = used(name)
	}
	return name
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := hasherPool.Get().(KeccakState)
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	hasherPool.Put(d)
	return b
}

var hasherPool = sync.Pool{
	New: func() any {
		return sha3.NewLegacyKeccak256().(KeccakState)
	},
}

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	d := hasherPool.Get().(KeccakState)
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	hasherPool.Put(d)
	return h
}
