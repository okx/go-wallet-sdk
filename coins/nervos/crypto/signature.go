package crypto

import (
	"encoding/hex"
	"math/big"
)

type SignatureData struct {
	V *big.Int
	R *big.Int
	S *big.Int

	ByteV byte
	ByteR []byte
	ByteS []byte
}

func (sd *SignatureData) ToHex() string {
	return hex.EncodeToString(sd.ToBytes())
}

func (sd SignatureData) ToBytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, sd.ByteR...)
	bytes = append(bytes, sd.ByteS...)

	// https://github.com/ethereum/go-ethereum/blob/master/crypto/signature_nocgo.go#L89
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sd.ByteV - 27
	bytes = append(bytes, v)
	return bytes
}
