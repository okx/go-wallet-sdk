package ed25519

import (
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"encoding/hex"
)

func GenerateKey() (ed25519.PrivateKey, error) {
	_, prv, err := ed25519.GenerateKey(crypto_rand.Reader)
	return prv, err
}

func PrivateKeyFromSeed(seedHex string) (ed25519.PrivateKey, error) {
	seedBytes, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, err
	}
	return ed25519.NewKeyFromSeed(seedBytes), nil
}

func PublicKeyFromSeed(seedHex string) (ed25519.PublicKey, error) {
	seedBytes, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, err
	}
	privateKey := ed25519.NewKeyFromSeed(seedBytes)
	return privateKey.Public().(ed25519.PublicKey), nil
}

func Sign(seedHex string, message []byte) (signature []byte, err error) {
	privateKey, err := PrivateKeyFromSeed(seedHex)
	if err != nil {
		return nil, err
	}
	return privateKey.Sign(crypto_rand.Reader, message, crypto.Hash(0))
}
