package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	btcec_ecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"math/big"
)

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes                     = wordBits / 8
	PubKeyBytesLenCompressed      = 33
	pubkeyCompressed         byte = 0x2 // y_bit + x coord
	pubkeyUncompressed       byte = 0x4 // x coord + y coord

)

var (
	secp256k1N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
)

type Secp256k1Key struct {
	PrivateKey *ecdsa.PrivateKey
}

func ZeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}

func SignAsRecoverable(value []byte, prvKey *btcec.PrivateKey) *SignatureData {
	sig, err := btcec_ecdsa.SignCompact(prvKey, value, false)
	if err != nil {
		return nil
	}
	V := sig[0]
	R := sig[1:33]
	S := sig[33:65]
	return &SignatureData{
		V:     new(big.Int).SetBytes([]byte{V}),
		R:     new(big.Int).SetBytes(R),
		S:     new(big.Int).SetBytes(S),
		ByteV: V,
		ByteR: R,
		ByteS: S,
	}
}

// Sign creates a recoverable ECDSA signature.
// The produced signature is in the 65-byte [R || S || V] format where V is 0 or 1.
func (k *Secp256k1Key) Sign(data []byte) ([]byte, error) {
	// da := SignAsRecoverable(data, (*btcec.PrivateKey)(k.PrivateKey))
	priKey, _ := btcec.PrivKeyFromBytes(k.PrivateKey.D.Bytes())
	da := SignAsRecoverable(data, priKey)
	return da.ToBytes(), nil
}

func HexToKey(hexkey string) (*Secp256k1Key, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToKey(b)
}

func ToKey(d []byte) (*Secp256k1Key, error) {
	return toKey(d, true)
}

func toKey(d []byte, strict bool) (*Secp256k1Key, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = secp256k1.S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, errors.New("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, errors.New("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return &Secp256k1Key{PrivateKey: priv}, nil
}

func RandomNew() (*Secp256k1Key, error) {
	randBytes := make([]byte, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		return nil, errors.New("key generation: could not read from random source: " + err.Error())
	}
	reader := bytes.NewReader(randBytes)
	priv, err := ecdsa.GenerateKey(secp256k1.S256(), reader)
	if err != nil {
		return nil, errors.New("key generation: ecdsa.GenerateKey failed: " + err.Error())
	}

	return &Secp256k1Key{PrivateKey: priv}, nil
}

func (k *Secp256k1Key) PubKey() []byte {
	pub := &k.PrivateKey.PublicKey
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}

	b := make([]byte, 0, PubKeyBytesLenCompressed)
	format := pubkeyCompressed
	if isOdd(pub.Y) {
		format |= 0x1
	}
	b = append(b, format)
	return paddedAppend(32, b, pub.X.Bytes())
}

func (k *Secp256k1Key) Bytes() []byte {
	return PaddedBigBytes(k.PrivateKey.D, k.PrivateKey.Params().BitSize/8)
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// isOdd returns whether the given integer is odd.
func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}
