package oasis

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/crypto/cbor"
)

func NewAddress(privateKeyHex string) (string, error) {
	bytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	privateKey := ed25519.PrivateKey(bytes)

	ctxData := append([]byte("oasis-core/address: staking"), uint8(0))
	pubBytes := privateKey[32:]

	hasher := sha512.New512_256()
	_, _ = hasher.Write(ctxData)
	_, _ = hasher.Write(pubBytes)
	hash := hasher.Sum([]byte{})
	addressBytes := append([]byte{uint8(0)}, hash[:20]...)

	converted, _ := bech32.ConvertBits(addressBytes, 8, 5, true)
	address, err := bech32.Encode("oasis", converted)

	if err != nil {
		return "", err
	}

	return address, nil
}

func SignTransaction(privateKeyHex, chainId string, tx *Transaction) *SignedTransaction {
	txBytes := cbor.Marshal(tx)

	ctx := fmt.Sprintf("oasis-core/consensus: tx for chain %s", chainId)
	ctxBytes := []byte(ctx)

	hasher := sha512.New512_256()
	_, _ = hasher.Write(ctxBytes)
	_, _ = hasher.Write(txBytes)
	hash := hasher.Sum([]byte{})

	privKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey := ed25519.PrivateKey(privKeyBytes)
	pubBytes := privateKey[32:]

	sig := ed25519.Sign(privateKey, hash)

	return &SignedTransaction{Signed{
		Blob: txBytes,
		Signature: Signature{
			PublicKey: pubBytes,
			Signature: sig,
		},
	}}
}
