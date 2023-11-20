package zil

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

const hrp = "zil"

func GetAddressFromPrivateKey(privKeyHex string) (string, error) {
	privBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return "", err
	}
	privateKey := secp256k1.PrivKeyFromBytes(privBytes)
	pubBytes := privateKey.PubKey().SerializeCompressed()
	hash := sha256.New()
	hash.Write(pubBytes)
	pubHash := hash.Sum(nil)[12:]
	conv, err := bech32.ConvertBits(pubHash, 8, 5, false)
	if err != nil {
		return "", err
	}
	address, err := bech32.Encode(hrp, conv)
	if err != nil {
		return "", err
	}
	return address, nil
}

func GetPublicKeyFromPrivateKey(privKeyHex string) (string, error) {
	privBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return "", err
	}
	privateKey := secp256k1.PrivKeyFromBytes(privBytes)
	pubBytes := privateKey.PubKey().SerializeCompressed()
	return hex.EncodeToString(pubBytes), nil
}

func FromBech32Addr(address string) (string, error) {
	h, data, err := bech32.Decode(address)
	if err != nil {
		return "", err
	}
	if h != hrp {
		return "", errors.New("expected hrp to be zil")
	}
	conv, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", err
	}
	buf := make([]byte, len(conv))
	for i := 0; i < len(conv); i++ {
		buf[i] = conv[i]
	}
	if len(buf) == 0 {
		return "", errors.New("could not convert buffer to bytes")
	}
	return hex.EncodeToString(buf), nil
}

func ToBech32Address(address string) (string, error) {
	if !IsAddress(address) {
		return "", errors.New("invalid address format")
	}
	data, _ := hex.DecodeString(address)
	conv, err := bech32.ConvertBits(data, 8, 5, false)
	if err != nil {
		return "", err
	}
	addr, err := bech32.Encode(hrp, conv)
	if err != nil {
		return "", err
	}

	return addr, nil
}
