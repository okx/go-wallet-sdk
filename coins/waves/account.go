package waves

import (
	"crypto/rand"
	"github.com/okx/go-wallet-sdk/coins/waves/crypto"
	"github.com/okx/go-wallet-sdk/coins/waves/types"
	"io"
)

// account
// https://docs.waves.tech/en/blockchain/account/address
const (
	// https://sourcegraph.com/github.com/wavesplatform/gowaves/-/blob/pkg/proto/addresses.go
	MainNetScheme   byte = 'W'
	TestNetScheme   byte = 'T'
	StageNetScheme  byte = 'S'
	CustomNetScheme byte = 'E'
)

// GenerateKeyPair generates a new key pair.
func GenerateKeyPair() (string, string, error) {
	seed := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		return "", "", err
	}

	secret, public, err := crypto.GenerateKeyPair(seed)
	if err != nil {
		return "", "", err
	}
	return secret.String(), public.String(), nil
}

// NewAddressFromString creates a WavesAddress from its string representation. This function checks that the address is valid.
func NewAddressFromString(s string) (types.WavesAddress, error) {
	return types.NewAddressFromString(s)
}

// NewAddressFromBytes creates a WavesAddress from the slice of bytes and checks that the result address is valid address.
func NewAddressFromBytes(b []byte) (types.WavesAddress, error) {
	return types.NewAddressFromBytes(b)
}

// GetAddress returns the String WavesAddress from the public key.
// The scheme is one of the MainNetScheme, TestNetScheme, StageNetScheme or CustomNetScheme.
// The public key is the base58 encoded string of the public key.
func GetAddress(scheme byte, pubKeyHash []byte) (string, error) {
	addr, err := types.NewAddressFromPublicKeyHash(scheme, pubKeyHash)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func ValidAddress(addr string) (bool, error) {
	wavesAddr, err := NewAddressFromString(addr)
	if err != nil {
		return false, err
	}
	return wavesAddr.Valid()
}
