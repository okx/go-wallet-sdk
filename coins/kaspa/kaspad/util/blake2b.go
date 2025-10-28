// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package util

import (
	"golang.org/x/crypto/blake2b"
)

// HashBlake2b calculates the hash blake2b(b).
func HashBlake2b(buf []byte) []byte {
	hashedBuf := blake2b.Sum256(buf)
	return hashedBuf[:]
}
