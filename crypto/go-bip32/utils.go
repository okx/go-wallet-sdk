package bip32

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/okx/go-wallet-sdk/crypto/btcd/btcec"
	"github.com/okx/go-wallet-sdk/crypto/ripemd160"
)

var (
	curve       = btcec.S256()
	curveParams = curve.Params()
)

//
// Hashes
//

func hashSha256(data []byte) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func hashDoubleSha256(data []byte) ([]byte, error) {
	hash1, err := hashSha256(data)
	if err != nil {
		return nil, err
	}

	hash2, err := hashSha256(hash1)
	if err != nil {
		return nil, err
	}
	return hash2, nil
}

func hashRipeMD160(data []byte) ([]byte, error) {
	hasher := ripemd160.New()
	_, err := io.WriteString(hasher, string(data))
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func hash160(data []byte) ([]byte, error) {
	hash1, err := hashSha256(data)
	if err != nil {
		return nil, err
	}

	hash2, err := hashRipeMD160(hash1)
	if err != nil {
		return nil, err
	}

	return hash2, nil
}

//
// Encoding
//

func checksum(data []byte) ([]byte, error) {
	hash, err := hashDoubleSha256(data)
	if err != nil {
		return nil, err
	}

	return hash[:4], nil
}

func addChecksumToBytes(data []byte) ([]byte, error) {
	checksum, err := checksum(data)
	if err != nil {
		return nil, err
	}
	return append(data, checksum...), nil
}

func base58Encode(data []byte) string {
	return base58.Encode(data)
}

func base58Decode(data string) []byte {
	return base58.Decode(data)
}

// Keys
func publicKeyForPrivateKey(key []byte) []byte {
	return compressPublicKey(curve.ScalarBaseMult(key))
}

func addPublicKeys(key1 []byte, key2 []byte) []byte {
	x1, y1 := expandPublicKey(key1)
	x2, y2 := expandPublicKey(key2)
	return compressPublicKey(curve.Add(x1, y1, x2, y2))
}

func addPrivateKeys(key1 []byte, key2 []byte) []byte {
	var key1Int big.Int
	var key2Int big.Int
	key1Int.SetBytes(key1)
	key2Int.SetBytes(key2)

	key1Int.Add(&key1Int, &key2Int)
	key1Int.Mod(&key1Int, curve.Params().N)

	b := key1Int.Bytes()
	if len(b) < 32 {
		extra := make([]byte, 32-len(b))
		b = append(extra, b...)
	}
	return b
}

func compressPublicKey(x *big.Int, y *big.Int) []byte {
	var key bytes.Buffer

	// Write header; 0x2 for even y value; 0x3 for odd
	key.WriteByte(byte(0x2) + byte(y.Bit(0)))

	// Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
	xBytes := x.Bytes()
	for i := 0; i < (PublicKeyCompressedLength - 1 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)

	return key.Bytes()
}

// As described at https://crypto.stackexchange.com/a/8916
func expandPublicKey(key []byte) (*big.Int, *big.Int) {
	Y := big.NewInt(0)
	X := big.NewInt(0)
	X.SetBytes(key[1:])

	// y^2 = x^3 + ax^2 + b
	// a = 0
	// => y^2 = x^3 + b
	ySquared := big.NewInt(0)
	ySquared.Exp(X, big.NewInt(3), nil)
	ySquared.Add(ySquared, curveParams.B)

	Y.ModSqrt(ySquared, curveParams.P)

	Ymod2 := big.NewInt(0)
	Ymod2.Mod(Y, big.NewInt(2))

	signY := uint64(key[0]) - 2
	if signY != Ymod2.Uint64() {
		Y.Sub(curveParams.P, Y)
	}

	return X, Y
}

func validatePrivateKey(key []byte) error {
	if fmt.Sprintf("%x", key) == "0000000000000000000000000000000000000000000000000000000000000000" || //if the key is zero
		bytes.Compare(key, curveParams.N.Bytes()) >= 0 || //or is outside of the curve
		len(key) != 32 { //or is too short
		return ErrInvalidPrivateKey
	}

	return nil
}

func validateChildPublicKey(key []byte) error {
	x, y := expandPublicKey(key)

	if x.Sign() == 0 || y.Sign() == 0 {
		return ErrInvalidPublicKey
	}

	return nil
}

//
// Numerical
//
func uint32Bytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)
	return bytes
}
