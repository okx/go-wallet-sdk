package ton

import (
	"crypto/ed25519"
	"errors"
	"strings"

	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
)

var (
	ErrInvalidMnemonic = errors.New("invalid mnemonic")
)

func NewAddress(seed []byte) (string, error) {
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	addr, err := wallet.AddressFromPubKey(pubKey, wallet.V4R2, wallet.DefaultSubwallet)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func AddressStrings(a string) ([]*address.AddrWithType, error) {
	addr, err := address.ParseAddr(a)
	if err != nil {
		return nil, err
	}
	return addr.Strings(), nil
}

func FromSeedV4R2(mnemonic, password string) (ed25519.PrivateKey, error) {
	if len(mnemonic) == 0 {
		return nil, ErrInvalidMnemonic
	}
	words := strings.Split(mnemonic, " ")
	if len(words) != 24 {
		return nil, ErrInvalidMnemonic
	}
	w, err := wallet.FromSeedWithPassword(words, password, wallet.V4R2)
	if err != nil {
		return nil, err
	}
	return w.PrivateKey(), nil
}

func ValidateAddress(addr string) bool {
	_, err := address.ParseAddr(addr)
	if err != nil {
		return false
	}
	return true
}

func VenomNewAddress(seed []byte) (string, error) {
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	addr, err := wallet.AddressFromPubKey(pubKey, wallet.VenomV3, wallet.VenomDefaultSubwallet)
	if err != nil {
		return "", err
	}

	return addr.RawString(), nil
}
