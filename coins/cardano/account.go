package cardano

import (
	"encoding/hex"

	"github.com/okx/go-wallet-sdk/coins/cardano/crypto"
	"github.com/okx/go-wallet-sdk/crypto/go-bip32"
	"github.com/okx/go-wallet-sdk/crypto/go-bip39"
)

func DerivePrvKey(mnemonic string, path string) (string, error) {
	entropy, err := bip39.EntropyFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	splitPath, err := bip32.ParseDerivationPath(path)
	if err != nil {
		return "", err
	}

	rootKey := crypto.NewXPrvKeyFromEntropy(entropy, "")
	accountKey := rootKey.
		Derive(splitPath[0]).
		Derive(splitPath[1]).
		Derive(splitPath[2])

	paymentKey := accountKey.
		Derive(0).
		Derive(splitPath[4])
	stakeKey := accountKey.
		Derive(2).
		Derive(splitPath[4])

	return paymentKey.PrvKey().String() + stakeKey.PrvKey().String(), nil
}

func NewAddressFromPrvKey(prvKeyHex string) (string, error) {
	prvKeyBytes, err := hex.DecodeString(prvKeyHex)
	if err != nil {
		return "", err
	}

	paymentPrvKey := crypto.PrvKey(prvKeyBytes[:64])
	payment, err := NewKeyCredential(paymentPrvKey.PubKey())
	if err != nil {
		return "", err
	}
	stakePrvKey := crypto.PrvKey(prvKeyBytes[64:])
	stake, err := NewKeyCredential(stakePrvKey.PubKey())
	if err != nil {
		return "", err
	}

	address, err := NewBaseAddress(Mainnet, payment, stake)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func NewAddressFromPubKey(pubKeyHex string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", err
	}

	paymentPubKey := crypto.PubKey(pubKeyBytes[0:32])
	payment, err := NewKeyCredential(paymentPubKey)
	if err != nil {
		return "", err
	}
	stakePubKey := crypto.PubKey(pubKeyBytes[32:])
	stake, err := NewKeyCredential(stakePubKey)
	if err != nil {
		return "", err
	}

	address, err := NewBaseAddress(Mainnet, payment, stake)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

func PubKeyFromPrvKey(prvKeyHex string) (string, error) {
	prvKeyBytes, err := hex.DecodeString(prvKeyHex)
	if err != nil {
		return "", err
	}

	paymentPrvKey := crypto.PrvKey(prvKeyBytes[:64])
	stakePrvKey := crypto.PrvKey(prvKeyBytes[64:])

	return paymentPrvKey.PubKey().String() + stakePrvKey.PubKey().String(), nil
}

func ValidateAddress(address string) bool {
	_, err := NewAddress(address)
	if err != nil {
		return false
	}
	return true
}
