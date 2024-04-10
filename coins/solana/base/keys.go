// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"

	"github.com/okx/go-wallet-sdk/crypto/base58"
)

type PrivateKey []byte

func MustPrivateKeyFromBase58(in string) PrivateKey {
	out, err := PrivateKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PrivateKeyFromBase58(privkey string) (PrivateKey, error) {
	res := base58.Decode(privkey)
	return res, nil
}

func (p PrivateKey) Bytes() []byte {
	return p[:]
}

func (k PrivateKey) String() string {
	return base58.Encode(k)
}

func NewRandomPrivateKey() (PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(crypto_rand.Reader)
	if err != nil {
		return nil, err
	}
	var publicKey PublicKey
	copy(publicKey[:], pub)
	return PrivateKey(priv), nil
}

func (k PrivateKey) Sign(payload []byte) (Signature, error) {
	p := ed25519.PrivateKey(k)
	signData, err := p.Sign(crypto_rand.Reader, payload, crypto.Hash(0))
	if err != nil {
		return Signature{}, err
	}

	var signature Signature
	copy(signature[:], signData)

	return signature, err
}

func (k PrivateKey) PublicKey() PublicKey {
	p := ed25519.PrivateKey(k)
	pub := p.Public().(ed25519.PublicKey)

	var publicKey PublicKey
	copy(publicKey[:], pub)

	return publicKey
}

type PublicKey [PublicKeyLength]byte

func PublicKeyFromBytes(in []byte) (out PublicKey) {
	byteCount := len(in)
	if byteCount == 0 {
		return
	}

	max := PublicKeyLength
	if byteCount < max {
		max = byteCount
	}

	copy(out[:], in[0:max])
	return
}

func MustPublicKeyFromBase58(in string) PublicKey {
	out, err := PublicKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PublicKeyFromBase58(in string) (out PublicKey, err error) {
	val := base58.Decode(in)
	if len(val) != PublicKeyLength {
		return out, fmt.Errorf("invalid length, expected %v, got %d", PublicKeyLength, len(val))
	}

	copy(out[:], val)
	return
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return []byte(base58.Encode(p[:])), nil
}

func (p *PublicKey) UnmarshalText(data []byte) (err error) {
	*p, err = PublicKeyFromBase58(string(data))
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", data, err)
	}
	return
}

func (p PublicKey) Equals(pb PublicKey) bool {
	return p == pb
}

// ToPointer returns a pointer to the pubkey.
func (p PublicKey) ToPointer() *PublicKey {
	return &p
}

func (p PublicKey) Bytes() []byte {
	return []byte(p[:])
}

var zeroPublicKey = PublicKey{}

// IsZero returns whether the public key is zero.
// NOTE: the System Program public key is also zero.
func (p PublicKey) IsZero() bool {
	return p == zeroPublicKey
}

func (p PublicKey) String() string {
	return base58.Encode(p[:])
}

// Short returns a shortened pubkey string,
// only including the first n chars, ellipsis, and the last n characters.
// NOTE: this is ONLY for visual representation for humans,
// and cannot be used for anything else.
func (p PublicKey) Short(n int) string {
	return formatShortPubkey(n, p)
}

func formatShortPubkey(n int, pubkey PublicKey) string {
	str := pubkey.String()
	if n > (len(str)/2)-1 {
		n = (len(str) / 2) - 1
	}
	if n < 2 {
		n = 2
	}
	return str[:n] + "..." + str[len(str)-n:]
}

type PublicKeySlice []PublicKey

// UniqueAppend appends the provided pubkey only if it is not
// already present in the slice.
// Returns true when the provided pubkey wasn't already present.
func (slice *PublicKeySlice) UniqueAppend(pubkey PublicKey) bool {
	if !slice.Has(pubkey) {
		slice.Append(pubkey)
		return true
	}
	return false
}

func (slice *PublicKeySlice) Append(pubkeys ...PublicKey) {
	*slice = append(*slice, pubkeys...)
}

func (slice PublicKeySlice) Has(pubkey PublicKey) bool {
	for _, key := range slice {
		if key.Equals(pubkey) {
			return true
		}
	}
	return false
}

// Split splits the slice into chunks of the specified size.
func (slice PublicKeySlice) Split(chunkSize int) []PublicKeySlice {
	divided := make([]PublicKeySlice, 0)
	if len(slice) == 0 || chunkSize < 1 {
		return divided
	}
	if len(slice) == 1 {
		return append(divided, slice)
	}

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		divided = append(divided, slice[i:end])
	}

	return divided
}

// GetAddedRemovedPubkeys accepts two slices of pubkeys (`previous` and `next`), and returns
// two slices:
// - `added` is the slice of pubkeys that are present in `next` but NOT present in `previous`.
// - `removed` is the slice of pubkeys that are present in `previous` but are NOT present in `next`.
func GetAddedRemovedPubkeys(previous PublicKeySlice, next PublicKeySlice) (added PublicKeySlice, removed PublicKeySlice) {
	added = make(PublicKeySlice, 0)
	removed = make(PublicKeySlice, 0)

	for _, prev := range previous {
		if !next.Has(prev) {
			removed = append(removed, prev)
		}
	}

	for _, nx := range next {
		if !previous.Has(nx) {
			added = append(added, nx)
		}
	}

	return
}

const (
	/// Number of bytes in a pubkey.
	PublicKeyLength = 32
	// Maximum length of derived pubkey seed.
	MaxSeedLength = 32
	// Maximum number of seeds.
	MaxSeeds = 16
	/// Number of bytes in a signature.
	SignatureLength = 64

	// // Maximum string length of a base58 encoded pubkey.
	// MaxBase58Length = 44
)

const PDA_MARKER = "ProgramDerivedAddress"

var ErrMaxSeedLengthExceeded = errors.New("Max seed length exceeded")

// Create a program address.
// Ported from https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L204
func CreateProgramAddress(seeds [][]byte, programID PublicKey) (PublicKey, error) {
	if len(seeds) > MaxSeeds {
		return PublicKey{}, ErrMaxSeedLengthExceeded
	}

	for _, seed := range seeds {
		if len(seed) > MaxSeedLength {
			return PublicKey{}, ErrMaxSeedLengthExceeded
		}
	}

	buf := []byte{}
	for _, seed := range seeds {
		buf = append(buf, seed...)
	}

	buf = append(buf, programID[:]...)
	buf = append(buf, []byte(PDA_MARKER)...)
	hash := sha256.Sum256(buf)

	if IsOnCurve(hash[:]) {
		return PublicKey{}, errors.New("invalid seeds; address must fall off the curve")
	}

	return PublicKeyFromBytes(hash[:]), nil
}

type incomparable [0]func()

// Point represents a point on the edwards25519 curve.
//
// This type works similarly to math/big.Int, and all arguments and receivers
// are allowed to alias.
//
// The zero value is NOT valid, and it may be used only as a receiver.
type Point struct {
	// The point is internally represented in extended coordinates (X, Y, Z, T)
	// where x = X/Z, y = Y/Z, and xy = T/Z per https://eprint.iacr.org/2008/522.
	x, y, z, t Element

	// Make the type not comparable (i.e. used with == or as a map key), as
	// equivalent points can be represented by different Go values.
	_ incomparable
}

// d is a constant in the curve equation.
var d, _ = new(Element).SetBytes([]byte{
	0xa3, 0x78, 0x59, 0x13, 0xca, 0x4d, 0xeb, 0x75,
	0xab, 0xd8, 0x41, 0x41, 0x4d, 0x0a, 0x70, 0x00,
	0x98, 0xe8, 0x79, 0x77, 0x79, 0x40, 0xc7, 0x8c,
	0x73, 0xfe, 0x6f, 0x2b, 0xee, 0x6c, 0x03, 0x52})

// SetBytes sets v = x, where x is a 32-byte encoding of v. If x does not
// represent a valid point on the curve, SetBytes returns nil and an error and
// the receiver is unchanged. Otherwise, SetBytes returns v.
//
// Note that SetBytes accepts all non-canonical encodings of valid points.
// That is, it follows decoding rules that match most implementations in
// the ecosystem rather than RFC 8032.
func (v *Point) SetBytes(x []byte) (*Point, error) {
	// Specifically, the non-canonical encodings that are accepted are
	//   1) the ones where the field element is not reduced (see the
	//      (*field.Element).SetBytes docs) and
	//   2) the ones where the x-coordinate is zero and the sign bit is set.
	//
	// This is consistent with crypto/ed25519/internal/edwards25519. Read more
	// at https://hdevalence.ca/blog/2020-10-04-its-25519am, specifically the
	// "Canonical A, R" section.

	y, err := new(Element).SetBytes(x)
	if err != nil {
		return nil, errors.New("edwards25519: invalid point encoding length")
	}

	// -x² + y² = 1 + dx²y²
	// x² + dx²y² = x²(dy² + 1) = y² - 1
	// x² = (y² - 1) / (dy² + 1)

	// u = y² - 1
	y2 := new(Element).Square(y)
	u := new(Element).Subtract(y2, feOne)

	// v = dy² + 1
	vv := new(Element).Multiply(y2, d)
	vv = vv.Add(vv, feOne)

	// x = +√(u/v)
	xx, wasSquare := new(Element).SqrtRatio(u, vv)
	if wasSquare == 0 {
		return nil, errors.New("edwards25519: invalid point encoding")
	}

	// Select the negative square root if the sign bit is set.
	xxNeg := new(Element).Negate(xx)
	xx = xx.Select(xxNeg, xx, int(x[31]>>7))

	v.x.Set(xx)
	v.y.Set(y)
	v.z.One()
	v.t.Multiply(xx, y) // xy = T / Z

	return v, nil
}

// Check if the provided `b` is on the ed25519 curve.
func IsOnCurve(b []byte) bool {
	_, err := new(Point).SetBytes(b)
	isOnCurve := err == nil
	return isOnCurve
}

// Find a valid program address and its corresponding bump seed.
func FindProgramAddress(seed [][]byte, programID PublicKey) (PublicKey, uint8, error) {
	var address PublicKey
	var err error
	bumpSeed := uint8(math.MaxUint8)
	for bumpSeed != 0 {
		address, err = CreateProgramAddress(append(seed, []byte{byte(bumpSeed)}), programID)
		if err == nil {
			return address, bumpSeed, nil
		}
		bumpSeed--
	}
	return PublicKey{}, bumpSeed, errors.New("unable to find a valid program address")
}

func FindAssociatedTokenAddress(wallet PublicKey, mint PublicKey, options ...string) (PublicKey, uint8, error) {
	return FindAssociatedTokenAddressAndBumpSeed(wallet, mint, SPLAssociatedTokenAccountProgramID, options...)
}

func FindAssociatedTokenAddressAndBumpSeed(walletAddress PublicKey, splTokenMintAddress PublicKey, programID PublicKey, options ...string) (PublicKey, uint8, error) {
	tokenProgramID := TokenProgramID
	if len(options) > 0 && options[0] == TOKEN2022 {
		tokenProgramID = Token2022ProgramID
	}
	return FindProgramAddress([][]byte{walletAddress[:], tokenProgramID[:], splTokenMintAddress[:]}, programID)
}
