package kava

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
)

const (
	HRP = "kava"
)

func NewAddress(privateKeyHex string) (string, error) {
	return cosmos.NewAddress(privateKeyHex, HRP)
}

func ValidateAddress(address string) bool {
	return cosmos.ValidateAddress(address, HRP)
}
