package near

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

func NewAccount() (address, seedHex string, err error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(publicKey), hex.EncodeToString(privateKey.Seed()), nil
}

// ValidateAddress NOTE:Address is not account id.
func ValidateAddress(addr string) bool {
	pubBytes, err := hex.DecodeString(addr)
	if err != nil {
		return false
	}
	if len(pubBytes) != 32 {
		return false
	}
	return true
}

func PrivateKeyToAddr(privateKey string) (string, error) {
	bytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	key := ed25519.PrivateKey(bytes)
	pubBytes := key[32:]
	return hex.EncodeToString(pubBytes), nil
}

func PrivateKeyToPublicKeyHex(privateKey string) (string, error) {
	privBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	pubBytes := privBytes[32:]
	return hex.EncodeToString(pubBytes), nil
}

func PublicKeyToAddress(publicKey string) (string, error) {
	publicKeyByte, err := hex.DecodeString(publicKey)
	if err != nil {
		return "", err
	}
	address := base58.Encode(publicKeyByte)
	return address, nil
}
