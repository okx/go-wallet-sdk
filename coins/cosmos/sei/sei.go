package sei

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
)

const (
	HRP = "sei"
)

func NewAddress(privateKeyHex string) (string, error) {
	return cosmos.NewAddress(privateKeyHex, HRP)
}

func ValidateAddress(address string) bool {
	return cosmos.ValidateAddress(address, HRP)
}
