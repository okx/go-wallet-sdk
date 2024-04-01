package near

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"regexp"
	"strings"
)

var (
	ErrInvalidPublicKey  = errors.New("invalid public key")
	ErrInvalidPrivateKey = errors.New("invalid private key")
	NearPrefix           = "ed25519:"
	Ed25519Prefix        = "ed25519"
)

func NewAccount() (address, pub, prvHex string, err error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", "", err
	}

	return hex.EncodeToString(publicKey), ExportPub(publicKey), ExportPrv(privateKey), nil
}

func ExportPrv(prv ed25519.PrivateKey) string {
	if prv == nil {
		return ""
	}
	return NearPrefix + base58.Encode(prv)
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func ExportPub(pub ed25519.PublicKey) string {
	if pub == nil {
		return ""
	}
	return NearPrefix + base58.Encode(pub)
}

var (
	regex, err = regexp.Compile("^(([a-z\\d]+[\\-_])*[a-z\\d]+\\.)*([a-z\\d]+[\\-_])*[a-z\\d]+$")
)

// ValidateAddress NOTE:Address is not account id.
func ValidateAddress(address string) bool {
	return len(address) >= 2 && len(address) <= 64 && regex.MatchString(address)
}

func PrivateKeyToAddr(privateKey string) (addr string, err error) {
	defer func() {
		if r := recover(); r != nil {
			addr, err = "", ErrInvalidPrivateKey
			return
		}
	}()
	if len(privateKey) == 0 {
		return "", ErrInvalidPrivateKey
	}
	if !strings.HasPrefix(privateKey, NearPrefix) {
		bytes, err := hex.DecodeString(privateKey)
		if err != nil {
			return "", err
		}

		key := ed25519.PrivateKey(bytes)
		pubBytes := key[32:]
		return hex.EncodeToString(pubBytes), nil
	}
	args := strings.Split(privateKey, ":")
	if len(args) != 2 || args[0] != Ed25519Prefix {
		return "", ErrInvalidPrivateKey
	}
	prv := base58.Decode(args[1])
	key := ed25519.PrivateKey(prv)
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, key[32:])
	return hex.EncodeToString(publicKey), nil
}

func PrivateKeyToPublicKey(privateKey string) (pk string, err error) {
	defer func() {
		if r := recover(); r != nil {
			pk, err = "", ErrInvalidPrivateKey
			return
		}
	}()
	if len(privateKey) == 0 {
		return "", ErrInvalidPrivateKey
	}
	if !strings.HasPrefix(privateKey, NearPrefix) {
		privBytes, err := hex.DecodeString(privateKey)
		if err != nil {
			return "", err
		}
		pubBytes := privBytes[32:]
		return NearPrefix + base58.Encode(pubBytes), nil
	}
	args := strings.Split(privateKey, ":")
	if len(args) != 2 || args[0] != Ed25519Prefix {
		return "", ErrInvalidPrivateKey
	}
	prv := base58.Decode(args[1])
	key := ed25519.PrivateKey(prv)
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, key[32:])
	return NearPrefix + base58.Encode(publicKey), nil
}

func PublicKeyToAddress(publicKey string) (addr string, err error) {
	defer func() {
		if r := recover(); r != nil {
			addr, err = "", ErrInvalidPublicKey
			return
		}
	}()
	if len(publicKey) == 0 {
		return "", ErrInvalidPublicKey
	}
	if !strings.HasPrefix(publicKey, NearPrefix) {
		publicKeyByte, err := hex.DecodeString(publicKey)
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(publicKeyByte), nil
	}
	args := strings.Split(publicKey, ":")
	if len(args) != 2 || args[0] != Ed25519Prefix {
		return "", ErrInvalidPublicKey
	}
	pk := base58.Decode(args[1])
	return hex.EncodeToString(pk), nil
}
