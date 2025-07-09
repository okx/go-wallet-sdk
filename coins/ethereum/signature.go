package ethereum

import (
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/crypto"
)

type SignatureData struct {
	V *big.Int
	R *big.Int
	S *big.Int

	ByteV byte
	ByteR []byte
	ByteS []byte
}

func NewSignatureData(msgHash []byte, publicKey string, r, s *big.Int) (*SignatureData, error) {
	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}

	pubKey, _ := btcec.ParsePubKey(pubBytes)
	sig, err := crypto.SignCompact(btcec.S256(), r, s, *pubKey, msgHash, false)
	if err != nil {
		return nil, err
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
	}, nil
}

func (sd *SignatureData) ToHex() string {
	return hex.EncodeToString(sd.ToBytes())
}

func (sd SignatureData) ToBytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, sd.ByteR...)
	bytes = append(bytes, sd.ByteS...)
	bytes = append(bytes, sd.ByteV)
	return bytes
}
