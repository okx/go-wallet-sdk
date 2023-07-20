package evmos

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"golang.org/x/crypto/sha3"
)

const (
	HRP = "evmos"
)

// The address generation method of eth is used
func NewAddress(privateKey string) (string, error) {
	pkBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	_, pb := btcec.PrivKeyFromBytes(pkBytes)
	pubBytes := pb.SerializeUncompressed()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes[1:])
	addressByte := hash.Sum(nil)
	address, err := bech32.EncodeFromBase256(HRP, addressByte[12:])
	if err != nil {
		return "", err
	}
	return address, nil
}

func ValidateAddress(address string) bool {
	return cosmos.ValidateAddress(address, HRP)
}
