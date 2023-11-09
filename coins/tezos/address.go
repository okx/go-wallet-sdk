package tezos

import (
	"github.com/okx/go-wallet-sdk/coins/tezos/types"
)

func ValidAddress(addr string) (bool, error) {
	tezosAddr, err := types.ParseAddress(addr)
	if err != nil {
		return false, err
	}
	return tezosAddr.IsValid(), nil
}

func GetAddress(publicKey string) (string, error) {
	return GetAddressByPublicKey(publicKey)
}

func GetAddressByPublicKey(publicKey string) (string, error) {
	key, err := types.ParseKey(publicKey)
	if err != nil {
		return "", err
	}
	return key.Address().String(), nil
}

func GetAddressByPrivateKey(privateKey string) (string, error) {
	key, err := types.ParsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	return key.Address().String(), nil
}

func GenerateKeyPair() (string, string, error) {
	key, err := types.GenerateKey(types.KeyTypeSecp256k1)
	if err != nil {
		return "", "", err
	}
	return key.String(), key.Public().String(), nil
}
