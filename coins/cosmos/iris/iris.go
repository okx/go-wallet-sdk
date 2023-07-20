package iris

import "github.com/okx/go-wallet-sdk/coins/cosmos"

const (
	HRP = "iaa"
)

func NewAddress(privateKeyHex string) (string, error) {
	return cosmos.NewAddress(privateKeyHex, HRP)
}

func ValidateAddress(address string) bool {
	return cosmos.ValidateAddress(address, HRP)
}
