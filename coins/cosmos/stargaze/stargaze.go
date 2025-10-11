package stargaze

import "github.com/emresenyuva/go-wallet-sdk/coins/cosmos"

const (
	HRP = "stars"
)

func NewAddress(privateKeyHex string) (string, error) {
	return cosmos.NewAddress(privateKeyHex, HRP, false)
}

func ValidateAddress(address string) bool {
	return cosmos.ValidateAddress(address, HRP)
}
