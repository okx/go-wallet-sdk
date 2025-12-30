package bip32

import (
	"errors"
	"strconv"
	"strings"

	"github.com/okx/go-wallet-sdk/crypto/btcd/v2/btcutil/hdkeychain"
)

func DerivePubKeyFromExtendedKey(extendedPublicKey string, path string) ([]byte, error) {
	extendedKey, err := hdkeychain.NewKeyFromString(extendedPublicKey)
	if err != nil {
		return nil, err
	}

	derivationPath, err := ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}

	for _, index := range derivationPath {
		if extendedKey, err = extendedKey.Derive(index); err != nil {
			return nil, err
		}
	}

	publicKey, err := extendedKey.ECPubKey()
	if err != nil {
		return nil, err
	}

	return publicKey.SerializeCompressed(), nil
}

func ParseDerivationPath(path string) ([]uint32, error) {
	var result []uint32

	components := strings.Split(path, "/")
	for _, component := range components {
		component = strings.TrimSpace(component)

		if component == "" || component == "m" {
			continue
		}

		index := strings.TrimSuffix(component, "'")
		value, err := strconv.ParseUint(index, 10, 32)
		if err != nil {
			return nil, err
		}

		if value >= hdkeychain.HardenedKeyStart {
			return nil, errors.New("child index too large")
		}

		if strings.HasSuffix(component, "'") {
			value += hdkeychain.HardenedKeyStart
		}

		result = append(result, uint32(value))
	}

	return result, nil
}
